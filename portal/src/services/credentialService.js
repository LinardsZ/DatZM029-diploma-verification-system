import api from '@/api';

export function postCredential(data) {
  return api().post('/credential', data);
}

export function verifyDiplomaHash(data) {
  return api().post('/verify/hash', data);
}

export function verifyDiplomaSignature(data) {
  return api().post('/verify/signature', data);
}

export function getCredentialList(id) {
  return api().get('/credentials', { params: { university: id } });
}

export function revokeCredential(id) {
  return api().patch(`/credential/${id}/revoke`);
}
