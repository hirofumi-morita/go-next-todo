'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/lib/auth-context';
import { todoApi, groupApi, Todo, Group } from '@/lib/api';
import Navbar from '@/components/Navbar';

export default function DashboardPage() {
  const { user, isLoading: authLoading } = useAuth();
  const router = useRouter();
  const [todos, setTodos] = useState<Todo[]>([]);
  const [groups, setGroups] = useState<Group[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [newTitle, setNewTitle] = useState('');
  const [newDescription, setNewDescription] = useState('');
  const [newGroupId, setNewGroupId] = useState<number | undefined>(undefined);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editTitle, setEditTitle] = useState('');
  const [editDescription, setEditDescription] = useState('');
  const [editGroupId, setEditGroupId] = useState<number | null>(null);
  const [error, setError] = useState('');
  const [filterGroup, setFilterGroup] = useState<number | 'all' | 'ungrouped'>('all');
  const [showGroupManager, setShowGroupManager] = useState(false);
  const [newGroupName, setNewGroupName] = useState('');
  const [newGroupDescription, setNewGroupDescription] = useState('');
  const [newGroupColor, setNewGroupColor] = useState('#3B82F6');
  const [editingGroupId, setEditingGroupId] = useState<number | null>(null);
  const [editGroupName, setEditGroupName] = useState('');
  const [editGroupDescription, setEditGroupDescription] = useState('');
  const [editGroupColor, setEditGroupColor] = useState('');

  useEffect(() => {
    if (!authLoading && !user) {
      router.push('/login');
    }
  }, [user, authLoading, router]);

  useEffect(() => {
    if (user) {
      fetchData();
    }
  }, [user]);

  const fetchData = async () => {
    const [todosRes, groupsRes] = await Promise.all([
      todoApi.getAll(),
      groupApi.getAll(),
    ]);
    if (todosRes.data) {
      setTodos(todosRes.data.todos || []);
    }
    if (groupsRes.data) {
      setGroups(groupsRes.data.groups || []);
    }
    if (todosRes.error) setError(todosRes.error);
    if (groupsRes.error) setError(groupsRes.error);
    setIsLoading(false);
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newTitle.trim()) return;

    const { data, error } = await todoApi.create(newTitle, newDescription, newGroupId);
    if (data) {
      setTodos([data.todo, ...todos]);
      setNewTitle('');
      setNewDescription('');
      setNewGroupId(undefined);
    }
    if (error) {
      setError(error);
    }
  };

  const handleToggleComplete = async (todo: Todo) => {
    const { data, error } = await todoApi.update(todo.id, {
      completed: !todo.completed,
    });
    if (data) {
      setTodos(todos.map((t) => (t.id === todo.id ? data.todo : t)));
    }
    if (error) {
      setError(error);
    }
  };

  const handleDelete = async (id: number) => {
    const { error } = await todoApi.delete(id);
    if (!error) {
      setTodos(todos.filter((t) => t.id !== id));
    } else {
      setError(error);
    }
  };

  const handleEdit = (todo: Todo) => {
    setEditingId(todo.id);
    setEditTitle(todo.title);
    setEditDescription(todo.description);
    setEditGroupId(todo.group_id ?? null);
  };

  const handleUpdate = async (id: number) => {
    const currentTodo = todos.find((t) => t.id === id);
    const currentGroupId = currentTodo?.group_id ?? null;
    const groupChanged = editGroupId !== currentGroupId;
    
    const updateData: { title: string; description: string; group_id?: string | null } = {
      title: editTitle,
      description: editDescription,
    };
    
    if (groupChanged) {
      updateData.group_id = editGroupId === null ? '' : String(editGroupId);
    }
    
    const { data, error } = await todoApi.update(id, updateData);
    if (data) {
      setTodos(todos.map((t) => (t.id === id ? data.todo : t)));
      setEditingId(null);
    }
    if (error) {
      setError(error);
    }
  };

  const handleCreateGroup = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newGroupName.trim()) return;

    const { data, error } = await groupApi.create(newGroupName, newGroupDescription, newGroupColor);
    if (data) {
      setGroups([...groups, data.group]);
      setNewGroupName('');
      setNewGroupDescription('');
      setNewGroupColor('#3B82F6');
    }
    if (error) {
      setError(error);
    }
  };

  const handleEditGroup = (group: Group) => {
    setEditingGroupId(group.id);
    setEditGroupName(group.name);
    setEditGroupDescription(group.description);
    setEditGroupColor(group.color);
  };

  const handleUpdateGroup = async (id: number) => {
    const { data, error } = await groupApi.update(id, {
      name: editGroupName,
      description: editGroupDescription,
      color: editGroupColor,
    });
    if (data) {
      setGroups(groups.map((g) => (g.id === id ? data.group : g)));
      setEditingGroupId(null);
    }
    if (error) {
      setError(error);
    }
  };

  const handleDeleteGroup = async (id: number) => {
    const { error } = await groupApi.delete(id);
    if (!error) {
      setGroups(groups.filter((g) => g.id !== id));
      setTodos(todos.map((t) => (t.group_id === id ? { ...t, group_id: null } : t)));
    } else {
      setError(error);
    }
  };

  const filteredTodos = todos.filter((todo) => {
    if (filterGroup === 'all') return true;
    if (filterGroup === 'ungrouped') return !todo.group_id;
    return todo.group_id === filterGroup;
  });

  const getGroupColor = (groupId: number | null | undefined) => {
    if (!groupId) return undefined;
    const group = groups.find((g) => g.id === groupId);
    return group?.color;
  };

  const getGroupName = (groupId: number | null | undefined) => {
    if (!groupId) return null;
    const group = groups.find((g) => g.id === groupId);
    return group?.name;
  };

  if (authLoading || isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <Navbar />
      <main className="max-w-4xl mx-auto py-8 px-4">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900">My TODOs</h1>
          <button
            onClick={() => setShowGroupManager(!showGroupManager)}
            className="bg-purple-600 text-white py-2 px-4 rounded-md hover:bg-purple-700"
          >
            {showGroupManager ? 'Hide Groups' : 'Manage Groups'}
          </button>
        </div>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {error}
            <button onClick={() => setError('')} className="float-right">&times;</button>
          </div>
        )}

        {showGroupManager && (
          <div className="bg-white rounded-lg shadow p-6 mb-8">
            <h2 className="text-lg font-semibold mb-4">Groups</h2>
            <form onSubmit={handleCreateGroup} className="flex gap-2 mb-4">
              <input
                type="text"
                placeholder="Group name"
                value={newGroupName}
                onChange={(e) => setNewGroupName(e.target.value)}
                className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
                required
              />
              <input
                type="text"
                placeholder="Description"
                value={newGroupDescription}
                onChange={(e) => setNewGroupDescription(e.target.value)}
                className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
              />
              <input
                type="color"
                value={newGroupColor}
                onChange={(e) => setNewGroupColor(e.target.value)}
                className="w-12 h-10 border border-gray-300 rounded-md cursor-pointer"
              />
              <button
                type="submit"
                className="bg-purple-600 text-white py-2 px-4 rounded-md hover:bg-purple-700"
              >
                Add
              </button>
            </form>
            <div className="space-y-2">
              {groups.map((group) => (
                <div key={group.id} className="flex items-center gap-2 p-2 bg-gray-50 rounded">
                  {editingGroupId === group.id ? (
                    <>
                      <input
                        type="text"
                        value={editGroupName}
                        onChange={(e) => setEditGroupName(e.target.value)}
                        className="flex-1 px-2 py-1 border border-gray-300 rounded"
                      />
                      <input
                        type="text"
                        value={editGroupDescription}
                        onChange={(e) => setEditGroupDescription(e.target.value)}
                        className="flex-1 px-2 py-1 border border-gray-300 rounded"
                      />
                      <input
                        type="color"
                        value={editGroupColor}
                        onChange={(e) => setEditGroupColor(e.target.value)}
                        className="w-10 h-8 border border-gray-300 rounded cursor-pointer"
                      />
                      <button
                        onClick={() => handleUpdateGroup(group.id)}
                        className="text-green-600 hover:text-green-800"
                      >
                        Save
                      </button>
                      <button
                        onClick={() => setEditingGroupId(null)}
                        className="text-gray-600 hover:text-gray-800"
                      >
                        Cancel
                      </button>
                    </>
                  ) : (
                    <>
                      <div
                        className="w-4 h-4 rounded-full"
                        style={{ backgroundColor: group.color }}
                      />
                      <span className="flex-1 font-medium">{group.name}</span>
                      <span className="text-gray-500 text-sm">{group.description}</span>
                      <button
                        onClick={() => handleEditGroup(group)}
                        className="text-blue-500 hover:text-blue-700 text-sm"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => handleDeleteGroup(group.id)}
                        className="text-red-500 hover:text-red-700 text-sm"
                      >
                        Delete
                      </button>
                    </>
                  )}
                </div>
              ))}
              {groups.length === 0 && (
                <p className="text-gray-500 text-center py-2">No groups yet</p>
              )}
            </div>
          </div>
        )}

        <form onSubmit={handleCreate} className="bg-white rounded-lg shadow p-6 mb-8">
          <h2 className="text-lg font-semibold mb-4">Add New TODO</h2>
          <div className="space-y-4">
            <input
              type="text"
              placeholder="Title"
              value={newTitle}
              onChange={(e) => setNewTitle(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
            <textarea
              placeholder="Description (optional)"
              value={newDescription}
              onChange={(e) => setNewDescription(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              rows={2}
            />
            <select
              value={newGroupId ?? ''}
              onChange={(e) => setNewGroupId(e.target.value ? Number(e.target.value) : undefined)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">No Group</option>
              {groups.map((group) => (
                <option key={group.id} value={group.id}>
                  {group.name}
                </option>
              ))}
            </select>
            <button
              type="submit"
              className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700"
            >
              Add TODO
            </button>
          </div>
        </form>

        <div className="mb-4 flex gap-2 flex-wrap">
          <button
            onClick={() => setFilterGroup('all')}
            className={`px-3 py-1 rounded-full text-sm ${
              filterGroup === 'all'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            All
          </button>
          <button
            onClick={() => setFilterGroup('ungrouped')}
            className={`px-3 py-1 rounded-full text-sm ${
              filterGroup === 'ungrouped'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            Ungrouped
          </button>
          {groups.map((group) => (
            <button
              key={group.id}
              onClick={() => setFilterGroup(group.id)}
              className={`px-3 py-1 rounded-full text-sm flex items-center gap-1 ${
                filterGroup === group.id
                  ? 'text-white'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              }`}
              style={filterGroup === group.id ? { backgroundColor: group.color } : undefined}
            >
              <div
                className="w-2 h-2 rounded-full"
                style={{ backgroundColor: group.color }}
              />
              {group.name}
            </button>
          ))}
        </div>

        <div className="space-y-4">
          {filteredTodos.length === 0 ? (
            <p className="text-gray-500 text-center py-8">No TODOs yet. Add one above!</p>
          ) : (
            filteredTodos.map((todo) => (
              <div
                key={todo.id}
                className={`bg-white rounded-lg shadow p-4 ${
                  todo.completed ? 'opacity-60' : ''
                }`}
                style={{
                  borderLeft: todo.group_id ? `4px solid ${getGroupColor(todo.group_id)}` : undefined,
                }}
              >
                {editingId === todo.id ? (
                  <div className="space-y-3">
                    <input
                      type="text"
                      value={editTitle}
                      onChange={(e) => setEditTitle(e.target.value)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md"
                    />
                    <textarea
                      value={editDescription}
                      onChange={(e) => setEditDescription(e.target.value)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md"
                      rows={2}
                    />
                    <select
                      value={editGroupId ?? ''}
                      onChange={(e) => setEditGroupId(e.target.value ? Number(e.target.value) : null)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md"
                    >
                      <option value="">No Group</option>
                      {groups.map((group) => (
                        <option key={group.id} value={group.id}>
                          {group.name}
                        </option>
                      ))}
                    </select>
                    <div className="flex gap-2">
                      <button
                        onClick={() => handleUpdate(todo.id)}
                        className="bg-green-500 text-white px-4 py-1 rounded-md hover:bg-green-600"
                      >
                        Save
                      </button>
                      <button
                        onClick={() => setEditingId(null)}
                        className="bg-gray-500 text-white px-4 py-1 rounded-md hover:bg-gray-600"
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                ) : (
                  <div className="flex items-start justify-between">
                    <div className="flex items-start gap-3">
                      <input
                        type="checkbox"
                        checked={todo.completed}
                        onChange={() => handleToggleComplete(todo)}
                        className="mt-1 h-5 w-5 rounded border-gray-300"
                      />
                      <div>
                        <h3
                          className={`font-medium ${
                            todo.completed ? 'line-through text-gray-500' : 'text-gray-900'
                          }`}
                        >
                          {todo.title}
                        </h3>
                        {todo.description && (
                          <p className="text-gray-600 text-sm mt-1">{todo.description}</p>
                        )}
                        <div className="flex items-center gap-2 mt-2">
                          {getGroupName(todo.group_id) && (
                            <span
                              className="text-xs px-2 py-0.5 rounded-full text-white"
                              style={{ backgroundColor: getGroupColor(todo.group_id) }}
                            >
                              {getGroupName(todo.group_id)}
                            </span>
                          )}
                          <p className="text-gray-400 text-xs">
                            Created: {new Date(todo.created_at).toLocaleDateString()}
                          </p>
                        </div>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <button
                        onClick={() => handleEdit(todo)}
                        className="text-blue-500 hover:text-blue-700"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => handleDelete(todo.id)}
                        className="text-red-500 hover:text-red-700"
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                )}
              </div>
            ))
          )}
        </div>
      </main>
    </div>
  );
}
