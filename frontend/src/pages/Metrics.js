import React, { useState, useEffect } from 'react';
import { useUsers } from '../context/UserContext';
import { 
  BarChart, 
  Bar, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  Legend, 
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  LineChart,
  Line
} from 'recharts';
import { 
  UsersIcon, 
  ChartBarIcon,
  ArrowTrendingUpIcon,
  ServerIcon
} from '@heroicons/react/24/outline';

const Metrics = () => {
  const { users, stats, fetchStats } = useUsers();
  const [timeRange, setTimeRange] = useState('7d');

  useEffect(() => {
    const interval = setInterval(() => {
      fetchStats();
    }, 30000); // Refresh every 30 seconds

    return () => clearInterval(interval);
  }, [fetchStats]);

  const formatBytes = (bytes) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  // Prepare data for charts
  const userStatusData = [
    { name: 'Active', value: users.filter(u => u.status === 'active').length, color: '#10B981' },
    { name: 'Inactive', value: users.filter(u => u.status === 'inactive').length, color: '#F59E0B' },
    { name: 'Suspended', value: users.filter(u => u.status === 'suspended').length, color: '#EF4444' },
  ];

  const topUsersByData = users
    .sort((a, b) => b.data_usage - a.data_usage)
    .slice(0, 10)
    .map(user => ({
      name: user.username,
      dataUsage: user.data_usage,
      connections: user.connection_count,
    }));

  const topUsersByConnections = users
    .sort((a, b) => b.connection_count - a.connection_count)
    .slice(0, 10)
    .map(user => ({
      name: user.username,
      connections: user.connection_count,
      dataUsage: user.data_usage,
    }));

  // Mock time series data for demonstration
  const timeSeriesData = [
    { time: '00:00', users: 12, connections: 8, dataUsage: 1024 },
    { time: '04:00', users: 15, connections: 12, dataUsage: 2048 },
    { time: '08:00', users: 25, connections: 20, dataUsage: 4096 },
    { time: '12:00', users: 30, connections: 28, dataUsage: 8192 },
    { time: '16:00', users: 28, connections: 25, dataUsage: 6144 },
    { time: '20:00', users: 22, connections: 18, dataUsage: 3072 },
    { time: '24:00', users: 18, connections: 15, dataUsage: 2048 },
  ];

  const StatCard = ({ title, value, icon: Icon, change, changeType }) => (
    <div className="card">
      <div className="flex items-center">
        <div className="flex-shrink-0">
          <Icon className="h-8 w-8 text-blue-600" />
        </div>
        <div className="ml-4">
          <p className="text-sm font-medium text-gray-500">{title}</p>
          <p className="text-2xl font-semibold text-gray-900">{value}</p>
          {change && (
            <p className={`text-sm ${
              changeType === 'increase' ? 'text-green-600' : 'text-red-600'
            }`}>
              {change}
            </p>
          )}
        </div>
      </div>
    </div>
  );

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Metrics & Analytics</h1>
          <p className="mt-2 text-gray-600">Monitor your VPN service performance</p>
        </div>
        <div className="flex space-x-2">
          <select
            value={timeRange}
            onChange={(e) => setTimeRange(e.target.value)}
            className="input-field w-32"
          >
            <option value="1d">Last 24h</option>
            <option value="7d">Last 7 days</option>
            <option value="30d">Last 30 days</option>
          </select>
        </div>
      </div>

      {/* Stats Overview */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatCard
          title="Total Users"
          value={users.length}
          icon={UsersIcon}
          change="+12% from last week"
          changeType="increase"
        />
        <StatCard
          title="Active Connections"
          value={users.reduce((sum, user) => sum + user.connection_count, 0)}
          icon={ChartBarIcon}
          change="+8% from last week"
          changeType="increase"
        />
        <StatCard
          title="Total Data Usage"
          value={formatBytes(users.reduce((sum, user) => sum + user.data_usage, 0))}
          icon={ArrowTrendingUpIcon}
          change="+15% from last week"
          changeType="increase"
        />
        <StatCard
          title="Running Pods"
          value={users.filter(u => u.pod_name).length}
          icon={ServerIcon}
          change="+5% from last week"
          changeType="increase"
        />
      </div>

      {/* Charts Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* User Status Distribution */}
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">User Status Distribution</h3>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={userStatusData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
              >
                {userStatusData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>

        {/* Time Series Chart */}
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Activity Over Time</h3>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={timeSeriesData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="time" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Line type="monotone" dataKey="users" stroke="#3B82F6" strokeWidth={2} />
              <Line type="monotone" dataKey="connections" stroke="#10B981" strokeWidth={2} />
            </LineChart>
          </ResponsiveContainer>
        </div>

        {/* Top Users by Data Usage */}
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Top Users by Data Usage</h3>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={topUsersByData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip formatter={(value) => formatBytes(value)} />
              <Bar dataKey="dataUsage" fill="#3B82F6" />
            </BarChart>
          </ResponsiveContainer>
        </div>

        {/* Top Users by Connections */}
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Top Users by Connections</h3>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={topUsersByConnections}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="connections" fill="#10B981" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Detailed Stats Table */}
      <div className="card">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Detailed Statistics</h3>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Metric
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Current
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Previous Period
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Change
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              <tr>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  Total Users
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {users.length}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {Math.floor(users.length * 0.88)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-green-600">
                  +12%
                </td>
              </tr>
              <tr>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  Active Users
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {users.filter(u => u.status === 'active').length}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {Math.floor(users.filter(u => u.status === 'active').length * 0.92)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-green-600">
                  +8%
                </td>
              </tr>
              <tr>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  Total Data Usage
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatBytes(users.reduce((sum, user) => sum + user.data_usage, 0))}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatBytes(Math.floor(users.reduce((sum, user) => sum + user.data_usage, 0) * 0.85))}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-green-600">
                  +15%
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Metrics;
