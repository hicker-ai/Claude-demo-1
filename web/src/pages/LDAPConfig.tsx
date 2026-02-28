import { useState, useEffect, useCallback } from 'react';
import {
  Card,
  Form,
  Input,
  InputNumber,
  Radio,
  Button,
  Badge,
  Descriptions,
  message,
  Typography,
  Space,
} from 'antd';
import { getLDAPConfig, updateLDAPConfig, getLDAPStatus } from '../api/ldap';
import type { LDAPConfig as LDAPConfigType, LDAPStatus } from '../api/ldap';

const { Title, Text } = Typography;

export default function LDAPConfig() {
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [config, setConfig] = useState<LDAPConfigType | null>(null);
  const [status, setStatus] = useState<LDAPStatus | null>(null);
  const [form] = Form.useForm<LDAPConfigType>();

  const fetchConfig = useCallback(async () => {
    setLoading(true);
    try {
      const resp = await getLDAPConfig();
      const data = resp.data?.data;
      setConfig(data);
      form.setFieldsValue(data);
    } catch {
      message.error('获取 LDAP 配置失败');
    } finally {
      setLoading(false);
    }
  }, [form]);

  const fetchStatus = useCallback(async () => {
    try {
      const resp = await getLDAPStatus();
      setStatus(resp.data?.data ?? null);
    } catch {
      // silently fail
    }
  }, []);

  useEffect(() => {
    void fetchConfig();
    void fetchStatus();
  }, [fetchConfig, fetchStatus]);

  const handleSave = async (values: LDAPConfigType) => {
    setSaving(true);
    try {
      await updateLDAPConfig(values);
      message.success('保存成功');
      void fetchConfig();
      void fetchStatus();
    } catch {
      message.error('保存失败');
    } finally {
      setSaving(false);
    }
  };

  const currentMode = Form.useWatch('mode', form) ?? config?.mode ?? 'OpenLDAP';
  const currentPort = Form.useWatch('port', form) ?? config?.port ?? 389;
  const currentBaseDN = Form.useWatch('base_dn', form) ?? config?.base_dn ?? 'dc=example,dc=com';

  const ldapSearchExample = currentMode === 'ActiveDirectory'
    ? `ldapsearch -H ldap://localhost:${currentPort} -D "CN=admin,${currentBaseDN}" -w password -b "${currentBaseDN}" "(sAMAccountName=username)"`
    : `ldapsearch -H ldap://localhost:${currentPort} -D "cn=admin,${currentBaseDN}" -w password -b "${currentBaseDN}" "(uid=username)"`;

  return (
    <div>
      <Title level={4} style={{ marginBottom: 16 }}>LDAP 配置</Title>

      <Space direction="vertical" size="middle" style={{ width: '100%' }}>
        <Card title="服务状态">
          <Descriptions column={3}>
            <Descriptions.Item label="运行状态">
              {status ? (
                <Badge
                  status={status.running ? 'success' : 'error'}
                  text={status.running ? '运行中' : '已停止'}
                />
              ) : (
                <Text type="secondary">加载中...</Text>
              )}
            </Descriptions.Item>
            <Descriptions.Item label="端口">
              {status?.port ?? '-'}
            </Descriptions.Item>
            <Descriptions.Item label="连接数">
              {status?.connections ?? '-'}
            </Descriptions.Item>
          </Descriptions>
        </Card>

        <Card title="配置" loading={loading}>
          <Form form={form} layout="vertical" onFinish={handleSave} style={{ maxWidth: 500 }}>
            <Form.Item
              name="base_dn"
              label="Base DN"
              rules={[{ required: true, message: '请输入 Base DN' }]}
            >
              <Input placeholder="dc=example,dc=com" />
            </Form.Item>
            <Form.Item
              name="mode"
              label="模式"
              rules={[{ required: true, message: '请选择模式' }]}
            >
              <Radio.Group>
                <Radio value="OpenLDAP">OpenLDAP</Radio>
                <Radio value="ActiveDirectory">Active Directory</Radio>
              </Radio.Group>
            </Form.Item>
            <Form.Item
              name="port"
              label="端口"
              rules={[{ required: true, message: '请输入端口号' }]}
            >
              <InputNumber min={1} max={65535} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={saving}>
                保存
              </Button>
            </Form.Item>
          </Form>
        </Card>

        <Card title="使用示例">
          <Text type="secondary" style={{ display: 'block', marginBottom: 8 }}>
            基于当前配置的 ldapsearch 命令示例：
          </Text>
          <div
            style={{
              background: '#f5f5f5',
              padding: 16,
              borderRadius: 8,
              fontFamily: 'monospace',
              fontSize: 13,
              wordBreak: 'break-all',
              whiteSpace: 'pre-wrap',
            }}
          >
            {ldapSearchExample}
          </div>
        </Card>
      </Space>
    </div>
  );
}
