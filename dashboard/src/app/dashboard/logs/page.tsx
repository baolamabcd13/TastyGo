'use client';

import { useState, useEffect } from 'react';
import axios from 'axios';

type ActivityLog = {
  id: number;
  user_id: number;
  username: string;
  activity_type: string;
  description: string;
  ip_address: string;
  created_at: string;
};

export default function LogsPage() {
  const [logs, setLogs] = useState<ActivityLog[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchLogs = async () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) return;

        const response = await axios.get('http://localhost:8080/api/admin/logs', {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });

        setLogs(response.data.data || []);
      } catch (err: any) {
        console.error('Error fetching logs:', err);
        setError(err.response?.data?.error || 'Không thể tải nhật ký hoạt động');
      } finally {
        setIsLoading(false);
      }
    };

    fetchLogs();
  }, []);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('vi-VN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    }).format(date);
  };

  const getActivityTypeLabel = (type: string) => {
    const types: Record<string, string> = {
      login: 'Đăng nhập',
      logout: 'Đăng xuất',
      create_user: 'Tạo người dùng',
      reset_password: 'Đặt lại mật khẩu',
      update_status: 'Cập nhật trạng thái',
      unlock_account: 'Mở khóa tài khoản'
    };
    
    return types[type] || type;
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-full">
        <span className="loading loading-spinner loading-lg text-primary"></span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="alert alert-error">
        <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span>{error}</span>
      </div>
    );
  }

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Nhật ký hoạt động</h1>
      
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <div className="overflow-x-auto">
            <table className="table table-zebra">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Người dùng</th>
                  <th>Hoạt động</th>
                  <th>Mô tả</th>
                  <th>Địa chỉ IP</th>
                  <th>Thời gian</th>
                </tr>
              </thead>
              <tbody>
                {logs.length > 0 ? (
                  logs.map((log) => (
                    <tr key={log.id}>
                      <td>{log.id}</td>
                      <td>{log.username}</td>
                      <td>
                        <div className="badge badge-primary text-white">
                          {getActivityTypeLabel(log.activity_type)}
                        </div>
                      </td>
                      <td>{log.description}</td>
                      <td>{log.ip_address}</td>
                      <td>{formatDate(log.created_at)}</td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={6} className="text-center py-4">
                      Không có dữ liệu
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}
