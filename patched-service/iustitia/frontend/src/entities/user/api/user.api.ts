import {
  api,
  ApiError,
  unwrap,
  type LoginRequest,
  type LoginResponse,
  type Role,
  type User,
  type UserListItem,
} from '@shared/api';

export const userApi = {
  login: (payload: LoginRequest): Promise<LoginResponse> =>
    unwrap(api.POST('/api/auth/login', { body: payload })),

  logout: async (): Promise<void> => {
    const result = await api.POST('/api/auth/logout');
    if (result.response.status === 204 || result.response.ok) {
      return;
    }
    throw new ApiError(result.response.status, 'Logout failed');
  },

  current: (): Promise<User> => unwrap(api.GET('/api/auth/me')),

  listByRole: (role: Role): Promise<UserListItem[]> =>
    unwrap(api.GET('/api/users', { params: { query: { role } } })),
};
