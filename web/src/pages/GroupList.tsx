import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Tree,
  Button,
  Modal,
  Form,
  Input,
  TreeSelect,
  message,
  Space,
  Typography,
  Card,
} from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { listGroups, createGroup } from '../api/groups';
import type { Group } from '../api/groups';
import type { DataNode } from 'antd/es/tree';

const { Title } = Typography;
const { TextArea } = Input;

function groupsToTreeData(groups: Group[]): DataNode[] {
  return groups.map((g) => ({
    key: g.id,
    title: g.name,
    children: g.children ? groupsToTreeData(g.children) : [],
  }));
}

interface TreeSelectNode {
  value: string;
  title: string;
  children?: TreeSelectNode[];
}

function groupsToTreeSelectData(groups: Group[]): TreeSelectNode[] {
  return groups.map((g) => ({
    value: g.id,
    title: g.name,
    children: g.children ? groupsToTreeSelectData(g.children) : [],
  }));
}

interface CreateGroupForm {
  name: string;
  description?: string;
  parent_id?: string;
}

export default function GroupList() {
  const navigate = useNavigate();
  const [groups, setGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [creating, setCreating] = useState(false);
  const [form] = Form.useForm<CreateGroupForm>();

  const fetchGroups = useCallback(async () => {
    setLoading(true);
    try {
      const resp = await listGroups();
      setGroups(resp.data?.data ?? []);
    } catch {
      message.error('获取用户组列表失败');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void fetchGroups();
  }, [fetchGroups]);

  const handleCreate = async (values: CreateGroupForm) => {
    setCreating(true);
    try {
      await createGroup(values);
      message.success('创建用户组成功');
      setModalOpen(false);
      form.resetFields();
      void fetchGroups();
    } catch {
      message.error('创建用户组失败');
    } finally {
      setCreating(false);
    }
  };

  const treeData = groupsToTreeData(groups);
  const treeSelectData = groupsToTreeSelectData(groups);

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>用户组管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
          创建用户组
        </Button>
      </div>

      <Card loading={loading}>
        {treeData.length > 0 ? (
          <Tree
            treeData={treeData}
            defaultExpandAll
            showLine
            onSelect={(selectedKeys) => {
              if (selectedKeys.length > 0) {
                navigate(`/groups/${selectedKeys[0]}`);
              }
            }}
          />
        ) : (
          <div style={{ textAlign: 'center', padding: 24, color: '#999' }}>
            暂无用户组，请点击"创建用户组"按钮创建
          </div>
        )}
      </Card>

      <Modal
        title="创建用户组"
        open={modalOpen}
        onCancel={() => {
          setModalOpen(false);
          form.resetFields();
        }}
        footer={null}
      >
        <Form form={form} layout="vertical" onFinish={handleCreate}>
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
          <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
            <Space>
              <Button onClick={() => { setModalOpen(false); form.resetFields(); }}>
                取消
              </Button>
              <Button type="primary" htmlType="submit" loading={creating}>
                创建
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
