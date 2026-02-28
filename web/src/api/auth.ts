import client from './client';

interface LoginResponse {
  code: number;
  message: string;
  data: {
    token: string;
  };
}

export async function login(username: string, password: string): Promise<LoginResponse> {
  const resp = await client.post<LoginResponse>('/auth/login', { username, password });
  return resp.data;
}

export async function logout(): Promise<void> {
  await client.post('/auth/logout');
}
