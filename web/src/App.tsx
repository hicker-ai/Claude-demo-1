import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import Login from './pages/Login';
import { isAuthenticated } from './store/auth';

function AuthGuard({ children }: { children: React.ReactNode }) {
  if (!isAuthenticated()) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
}

function Placeholder({ title }: { title: string }) {
  return <div style={{ padding: 24 }}><h2>{title}</h2><p>页面开发中...</p></div>;
}

export default function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/users" element={<AuthGuard><Placeholder title="用户管理" /></AuthGuard>} />
          <Route path="/groups" element={<AuthGuard><Placeholder title="用户组管理" /></AuthGuard>} />
          <Route path="/ldap-config" element={<AuthGuard><Placeholder title="LDAP 配置" /></AuthGuard>} />
          <Route path="/" element={<Navigate to="/users" replace />} />
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  );
}
