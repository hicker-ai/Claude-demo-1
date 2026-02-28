import client from './client';

export interface Group {
  id: string;
  name: string;
  description: string;
  parent_id: string | null;
  children?: Group[];
  created_at: string;
  updated_at: string;
}

export const listGroups = () => client.get('/groups');
export const getGroup = (id: string) => client.get(`/groups/${id}`);
export const createGroup = (data: { name: string; description?: string; parent_id?: string }) =>
  client.post('/groups', data);
export const updateGroup = (id: string, data: { name?: string; description?: string; parent_id?: string }) =>
  client.put(`/groups/${id}`, data);
export const deleteGroup = (id: string) => client.delete(`/groups/${id}`);
export const getGroupMembers = (id: string) => client.get(`/groups/${id}/members`);
export const addMembers = (id: string, userIds: string[]) =>
  client.post(`/groups/${id}/members`, { user_ids: userIds });
export const removeMember = (id: string, userId: string) =>
  client.delete(`/groups/${id}/members/${userId}`);
