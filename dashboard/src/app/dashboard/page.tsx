'use client';

import { useState, useEffect } from 'react';
import axios from 'axios';

export default function DashboardPage() {
  const [profile, setProfile] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) return;

        const response = await axios.get('http://localhost:8080/api/profile', {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });

        setProfile(response.data);
      } catch (err: any) {
        console.error('Error fetching profile:', err);
        setError('Không thể tải thông tin người dùng');
      } finally {
        setIsLoading(false);
      }
    };

    fetchProfile();
  }, []);

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
      <h1 className="text-3xl font-bold mb-6">Bảng điều khiển</h1>
      
      {profile && (
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title text-xl mb-4">Thông tin tài khoản</h2>
            
            <div className="overflow-x-auto">
              <table className="table">
                <tbody>
                  <tr>
                    <td className="font-medium">Tên người dùng:</td>
                    <td>{profile.username}</td>
                  </tr>
                  <tr>
                    <td className="font-medium">Email:</td>
                    <td>{profile.email}</td>
                  </tr>
                  <tr>
                    <td className="font-medium">Vai trò:</td>
                    <td>
                      <div className="badge badge-primary text-white">
                        {profile.role === 'superadmin' ? 'Super Admin' : 'Admin'}
                      </div>
                    </td>
                  </tr>
                  {profile.profile && profile.profile.full_name && (
                    <tr>
                      <td className="font-medium">Họ tên:</td>
                      <td>{profile.profile.full_name}</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mt-8">
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">Quản lý Admin</h2>
            <p>Quản lý tài khoản Admin trong hệ thống</p>
            <div className="card-actions justify-end">
              <button className="btn btn-primary text-white" onClick={() => window.location.href = '/dashboard/admins'}>
                Xem chi tiết
              </button>
            </div>
          </div>
        </div>
        
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">Nhật ký hoạt động</h2>
            <p>Xem lịch sử hoạt động của người dùng</p>
            <div className="card-actions justify-end">
              <button className="btn btn-primary text-white" onClick={() => window.location.href = '/dashboard/logs'}>
                Xem chi tiết
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
