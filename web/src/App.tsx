import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import Login from './pages/Login';
import Layout from './components/Layout';
import UserList from './pages/UserList';
import UserDetail from './pages/UserDetail';
import GroupList from './pages/GroupList';
import GroupDetail from './pages/GroupDetail';
import LDAPConfig from './pages/LDAPConfig';
import { isAuthenticated } from './store/auth';

function AuthGuard({ children }: { children: React.ReactNode }) {
  if (!isAuthenticated()) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
}

export default function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route
            element={
              <AuthGuard>
                <Layout />
              </AuthGuard>
            }
          >
            <Route path="/users" element={<UserList />} />
            <Route path="/users/:id" element={<UserDetail />} />
            <Route path="/groups" element={<GroupList />} />
            <Route path="/groups/:id" element={<GroupDetail />} />
            <Route path="/ldap-config" element={<LDAPConfig />} />
          </Route>
          <Route path="/" element={<Navigate to="/users" replace />} />
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  );
}
