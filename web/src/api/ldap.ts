import client from './client';

export interface LDAPConfig {
  base_dn: string;
  mode: string;
  port: number;
}

export interface LDAPStatus {
  running: boolean;
  port: number;
  connections: number;
}

export const getLDAPConfig = () => client.get('/ldap/config');
export const updateLDAPConfig = (data: { base_dn: string; mode: string; port: number }) =>
  client.put('/ldap/config', data);
export const getLDAPStatus = () => client.get('/ldap/status');
