import api from '@/api';

export function postCredential(data) {
  return api().post('/credential', data);
}
