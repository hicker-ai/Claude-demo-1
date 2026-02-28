import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Card,
  Form,
  Input,
  Button,
  TreeSelect,
  Table,
  Modal,
  Select,
  message,
  Space,
  Popconfirm,
  Typography,
  Divider,
} from 'antd';
import { ArrowLeftOutlined, PlusOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import {
  getGroup,
  updateGroup,
  listGroups,
  getGroupMembers,
  addMembers,
  removeMember,
} from '../api/groups';
import type { Group } from '../api/groups';
import { listUsers } from '../api/users';
import type { User } from '../api/users';

const { Title } = Typography;
const { TextArea } = Input;

interface GroupForm {
  name: string;
  description?: string;
  parent_id?: string;
}

interface TreeSelectNode {
  value: string;
  title: string;
  children?: TreeSelectNode[];
}

function groupsToTreeSelectData(groups: Group[], excludeId?: string): TreeSelectNode[] {
  return groups
    .filter((g) => g.id !== excludeId)
    .map((g) => ({
      value: g.id,
      title: g.name,
      children: g.children ? groupsToTreeSelectData(g.children, excludeId) : [],
    }));
}

export default function GroupDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [, setGroup] = useState<Group | null>(null);
  const [allGroups, setAllGroups] = useState<Group[]>([]);
  const [members, setMembers] = useState<User[]>([]);
  const [addModalOpen, setAddModalOpen] = useState(false);
  const [adding, setAdding] = useState(false);
  const [selectedUserIds, setSelectedUserIds] = useState<string[]>([]);
  const [userOptions, setUserOptions] = useState<{ label: string; value: string }[]>([]);
  const [form] = Form.useForm<GroupForm>();

  const fetchGroup = useCallback(async () => {
    if (!id) return;
    setLoading(true);
    try {
      const resp = await getGroup(id);
      const data = resp.data?.data;
      setGroup(data);
      form.setFieldsValue({
        name: data.name,
        description: data.description,
        parent_id: data.parent_id ?? undefined,
      });
    } catch {
      message.error('获取用户组信息失败');
    } finally {
      setLoading(false);
    }
  }, [id, form]);

  const fetchAllGroups = useCallback(async () => {
    try {
      const resp = await listGroups();
      setAllGroups(resp.data?.data ?? []);
    } catch {
      // silently fail
    }
  }, []);

  const fetchMembers = useCallback(async () => {
    if (!id) return;
    try {
      const resp = await getGroupMembers(id);
      setMembers(resp.data?.data ?? []);
    } catch {
      message.error('获取成员列表失败');
    }
  }, [id]);

  useEffect(() => {
    void fetchGroup();
    void fetchAllGroups();
    void fetchMembers();
  }, [fetchGroup, fetchAllGroups, fetchMembers]);

  const handleSave = async (values: GroupForm) => {
    if (!id) return;
    setSaving(true);
    try {
      await updateGroup(id, values);
      message.success('保存成功');
      void fetchGroup();
      void fetchAllGroups();
    } catch {
      message.error('保存失败');
    } finally {
      setSaving(false);
    }
  };

  const handleRemoveMember = async (userId: string) => {
    if (!id) return;
    try {
      await removeMember(id, userId);
      message.success('移除成功');
      void fetchMembers();
    } catch {
      message.error('移除失败');
    }
  };

  const handleSearchUsers = async (searchText: string) => {
    try {
      const resp = await listUsers({ page: 1, page_size: 20, search: searchText || undefined });
      const items: User[] = resp.data?.data?.items ?? [];
      setUserOptions(
        items.map((u) => ({
          label: `${u.username} (${u.display_name})`,
          value: u.id,
        }))
      );
    } catch {
      // silently fail
    }
  };

  const handleAddMembers = async () => {
    if (!id || selectedUserIds.length === 0) return;
    setAdding(true);
    try {
      await addMembers(id, selectedUserIds);
      message.success('添加成功');
      setAddModalOpen(false);
      setSelectedUserIds([]);
      void fetchMembers();
    } catch {
      message.error('添加失败');
    } finally {
      setAdding(false);
    }
  };

  const memberColumns: ColumnsType<User> = [
    { title: '用户名', dataIndex: 'username', key: 'username' },
    { title: '显示名称', dataIndex: 'display_name', key: 'display_name' },
    { title: '邮箱', dataIndex: 'email', key: 'email' },
    {
      title: '操作',
      key: 'actions',
      render: (_, record) => (
        <Popconfirm
          title="确认移除"
          description="确定要将该用户从用户组中移除吗？"
          onConfirm={() => handleRemoveMember(record.id)}
          okText="确认"
          cancelText="取消"
        >
          <Button type="link" size="small" danger>
            移除
          </Button>
        </Popconfirm>
      ),
    },
  ];

  const treeSelectData = groupsToTreeSelectData(allGroups, id);

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/groups')}>
          返回
        </Button>
        <Title level={4} style={{ margin: 0 }}>用户组详情</Title>
      </Space>

      <Card title="基本信息" loading={loading} style={{ marginBottom: 16 }}>
        <Form form={form} layout="vertical" onFinish={handleSave} style={{ maxWidth: 500 }}>
          <Form.Item
            name="name"
            label="组名称"
            rules={[{ required: true, message: '请输入组名称' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <TextArea rows={3} />
          </Form.Item>
          <Form.Item name="parent_id" label="上级用户组">
            <TreeSelect
              treeData={treeSelectData}
              placeholder="请选择上级用户组（可选）"
              allowClear
              treeDefaultExpandAll
            />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={saving}>
              保存
            </Button>
          </Form.Item>
        </Form>
      </Card>

      <Card
        title="成员列表"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={() => {
            setAddModalOpen(true);
            void handleSearchUsers('');
          }}>
            添加成员
          </Button>
        }
      >
        <Divider style={{ margin: '0 0 16px 0' }} />
        <Table
          rowKey="id"
          columns={memberColumns}
          dataSource={members}
          pagination={false}
        />
      </Card>

      <Modal
        title="添加成员"
        open={addModalOpen}
        onCancel={() => {
          setAddModalOpen(false);
          setSelectedUserIds([]);
        }}
        onOk={handleAddMembers}
        confirmLoading={adding}
        okText="添加"
        cancelText="取消"
      >
        <div style={{ marginBottom: 8 }}>选择要添加的用户：</div>
        <Select
          mode="multiple"
          style={{ width: '100%' }}
          placeholder="搜索并选择用户"
          value={selectedUserIds}
          onChange={setSelectedUserIds}
          onSearch={handleSearchUsers}
          options={userOptions}
          filterOption={false}
          showSearch
        />
      </Modal>
    </div>
  );
}
