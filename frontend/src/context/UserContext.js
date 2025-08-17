import React, { createContext, useContext, useReducer, useEffect } from 'react';
import axios from 'axios';
import toast from 'react-hot-toast';

const UserContext = createContext();

const initialState = {
  users: [],
  loading: false,
  error: null,
  stats: null,
};

const userReducer = (state, action) => {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, loading: action.payload };
    case 'SET_ERROR':
      return { ...state, error: action.payload, loading: false };
    case 'SET_USERS':
      return { ...state, users: action.payload, loading: false, error: null };
    case 'ADD_USER':
      return { ...state, users: [...state.users, action.payload] };
    case 'DELETE_USER':
      return { ...state, users: state.users.filter(user => user.id !== action.payload) };
    case 'SET_STATS':
      return { ...state, stats: action.payload };
    default:
      return state;
  }
};

export const UserProvider = ({ children }) => {
  const [state, dispatch] = useReducer(userReducer, initialState);

  const api = axios.create({
    baseURL: '/api/v1',
    timeout: 10000,
  });

  const fetchUsers = async () => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true });
      const response = await api.get('/users');
      dispatch({ type: 'SET_USERS', payload: response.data.users });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: error.message });
      toast.error('Failed to fetch users');
    }
  };

  const createUser = async (userData) => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true });
      const response = await api.post('/users', userData);
      dispatch({ type: 'ADD_USER', payload: response.data.user });
      toast.success('User created successfully');
      return response.data.user;
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: error.message });
      toast.error(error.response?.data?.error || 'Failed to create user');
      throw error;
    }
  };

  const deleteUser = async (userId) => {
    try {
      await api.delete(`/users/${userId}`);
      dispatch({ type: 'DELETE_USER', payload: userId });
      toast.success('User deleted successfully');
    } catch (error) {
      toast.error('Failed to delete user');
      throw error;
    }
  };

  const downloadConfig = async (userId, username) => {
    try {
      const response = await api.get(`/users/${userId}/config`, {
        responseType: 'blob',
      });
      
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `vpn-${username}.conf`);
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);
      
      toast.success('Configuration downloaded successfully');
    } catch (error) {
      toast.error('Failed to download configuration');
    }
  };

  const fetchStats = async () => {
    try {
      const response = await api.get('/stats');
      dispatch({ type: 'SET_STATS', payload: response.data.stats });
    } catch (error) {
      console.error('Failed to fetch stats:', error);
    }
  };

  useEffect(() => {
    fetchUsers();
    fetchStats();
  }, []);

  const value = {
    ...state,
    fetchUsers,
    createUser,
    deleteUser,
    downloadConfig,
    fetchStats,
  };

  return (
    <UserContext.Provider value={value}>
      {children}
    </UserContext.Provider>
  );
};

export const useUsers = () => {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error('useUsers must be used within a UserProvider');
  }
  return context;
};
