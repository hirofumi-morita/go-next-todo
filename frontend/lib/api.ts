const API_BASE = '/api';

interface ApiResponse<T> {
  data?: T;
  error?: string;
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
    ...options.headers,
  };

  try {
    const response = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers,
    });

    const data = await response.json();

    if (!response.ok) {
      return { error: data.error || 'An error occurred' };
    }

    return { data };
  } catch (error) {
    return { error: 'Network error occurred' };
  }
}

export const authApi = {
  register: (email: string, password: string) =>
    request<{ token: string; user: User; message: string }>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),

  login: (email: string, password: string) =>
    request<{ token: string; user: User; message: string }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),

  getMe: () => request<{ user: User }>('/me'),
};

export const todoApi = {
  getAll: () => request<{ todos: Todo[] }>('/todos'),

  get: (id: number) => request<{ todo: Todo }>(`/todos/${id}`),

  create: (title: string, description: string) =>
    request<{ todo: Todo }>('/todos', {
      method: 'POST',
      body: JSON.stringify({ title, description }),
    }),

  update: (id: number, data: Partial<Todo>) =>
    request<{ todo: Todo }>(`/todos/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  delete: (id: number) =>
    request<{ message: string }>(`/todos/${id}`, {
      method: 'DELETE',
    }),
};

export const adminApi = {
  getUsers: () => request<{ users: User[] }>('/admin/users'),

  getUser: (id: number) => request<{ user: User }>(`/admin/users/${id}`),

  deleteUser: (id: number) =>
    request<{ message: string }>(`/admin/users/${id}`, {
      method: 'DELETE',
    }),

  updateUserAdmin: (id: number, isAdmin: boolean) =>
    request<{ user: User }>(`/admin/users/${id}`, {
      method: 'PATCH',
      body: JSON.stringify({ is_admin: isAdmin }),
    }),
};

export interface User {
  id: number;
  email: string;
  is_admin: boolean;
  created_at?: string;
  todos?: Todo[];
}

export interface Todo {
  id: number;
  title: string;
  description: string;
  completed: boolean;
  user_id: number;
  created_at: string;
  updated_at: string;
}
