import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Card,
  Form,
  Input,
  Button,
  message,
  Space,
  Table,
  Typography,
  Divider,
} from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { getUser, updateUser, changePassword, getUserGroups } from '../api/users';
import type { User, UpdateUserReq } from '../api/users';
import type { Group } from '../api/groups';

const { Title } = Typography;

interface PasswordForm {
  old_password: string;
  new_password: string;
  confirm_password: string;
}

export default function UserDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [changingPwd, setChangingPwd] = useState(false);
  const [user, setUser] = useState<User | null>(null);
  const [groups, setGroups] = useState<Group[]>([]);
  const [infoForm] = Form.useForm<UpdateUserReq>();
  const [pwdForm] = Form.useForm<PasswordForm>();

  const fetchUser = useCallback(async () => {
    if (!id) return;
    setLoading(true);
    try {
      const resp = await getUser(id);
      const data = resp.data?.data;
      setUser(data);
      infoForm.setFieldsValue({
        display_name: data.display_name,
        email: data.email,
        phone: data.phone,
      });
    } catch {
      message.error('获取用户信息失败');
    } finally {
      setLoading(false);
    }
  }, [id, infoForm]);

  const fetchGroups = useCallback(async () => {
    if (!id) return;
    try {
      const resp = await getUserGroups(id);
      setGroups(resp.data?.data ?? []);
    } catch {
      message.error('获取用户组信息失败');
    }
  }, [id]);

  useEffect(() => {
    void fetchUser();
    void fetchGroups();
  }, [fetchUser, fetchGroups]);

  const handleSave = async (values: UpdateUserReq) => {
    if (!id) return;
    setSaving(true);
    try {
      await updateUser(id, values);
      message.success('保存成功');
      void fetchUser();
    } catch {
      message.error('保存失败');
    } finally {
      setSaving(false);
    }
  };

  const handleChangePassword = async (values: PasswordForm) => {
    if (!id) return;
    if (values.new_password !== values.confirm_password) {
      message.error('两次输入的密码不一致');
      return;
    }
    setChangingPwd(true);
    try {
      await changePassword(id, {
        old_password: values.old_password,
        new_password: values.new_password,
      });
      message.success('密码修改成功');
      pwdForm.resetFields();
    } catch {
      message.error('密码修改失败');
    } finally {
      setChangingPwd(false);
    }
  };

  const groupColumns: ColumnsType<Group> = [
    { title: '组名称', dataIndex: 'name', key: 'name' },
    { title: '描述', dataIndex: 'description', key: 'description' },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
    },
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/users')}>
          返回
        </Button>
        <Title level={4} style={{ margin: 0 }}>用户详情</Title>
      </Space>

      <Card title="基本信息" loading={loading} style={{ marginBottom: 16 }}>
        <Form form={infoForm} layout="vertical" onFinish={handleSave} style={{ maxWidth: 500 }}>
          <Form.Item label="用户名">
            <Input value={user?.username ?? ''} disabled />
          </Form.Item>
          <Form.Item
            name="display_name"
            label="显示名称"
            rules={[{ required: true, message: '请输入显示名称' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="email"
            label="邮箱"
            rules={[
              { required: true, message: '请输入邮箱' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input />
          </Form.Item>
          <Form.Item name="phone" label="手机号">
            <Input />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={saving}>
              保存
            </Button>
          </Form.Item>
        </Form>
      </Card>

      <Card title="修改密码" style={{ marginBottom: 16 }}>
        <Form form={pwdForm} layout="vertical" onFinish={handleChangePassword} style={{ maxWidth: 500 }}>
          <Form.Item
            name="old_password"
            label="原密码"
            rules={[{ required: true, message: '请输入原密码' }]}
          >
            <Input.Password />
          </Form.Item>
          <Form.Item
            name="new_password"
            label="新密码"
            rules={[{ required: true, message: '请输入新密码' }]}
          >
            <Input.Password />
          </Form.Item>
          <Form.Item
            name="confirm_password"
            label="确认新密码"
            rules={[
              { required: true, message: '请确认新密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('new_password') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'));
                },
              }),
            ]}
          >
            <Input.Password />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={changingPwd}>
              修改密码
            </Button>
          </Form.Item>
        </Form>
      </Card>

      <Card title="所属用户组">
        <Divider style={{ margin: '0 0 16px 0' }} />
        <Table
          rowKey="id"
          columns={groupColumns}
          dataSource={groups}
          pagination={false}
        />
      </Card>
    </div>
  );
}
