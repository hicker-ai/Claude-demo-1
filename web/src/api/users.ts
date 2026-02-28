import client from './client';

export interface User {
  id: string;
  username: string;
  display_name: string;
  email: string;
  phone: string;
  status: 'enabled' | 'disabled';
  created_at: string;
  updated_at: string;
}

export interface CreateUserReq {
  username: string;
  display_name: string;
  email: string;
  password: string;
  phone?: string;
}

export interface UpdateUserReq {
  display_name?: string;
  email?: string;
  phone?: string;
}

export const listUsers = (params: { page: number; page_size: number; search?: string }) =>
  client.get('/users', { params });

export const getUser = (id: string) => client.get(`/users/${id}`);
export const createUser = (data: CreateUserReq) => client.post('/users', data);
export const updateUser = (id: string, data: UpdateUserReq) => client.put(`/users/${id}`, data);
export const deleteUser = (id: string) => client.delete(`/users/${id}`);
export const changePassword = (id: string, data: { old_password: string; new_password: string }) =>
  client.put(`/users/${id}/password`, data);
export const setUserStatus = (id: string, status: 'enabled' | 'disabled') =>
  client.put(`/users/${id}/status`, { status });
export const getUserGroups = (id: string) => client.get(`/users/${id}/groups`);
