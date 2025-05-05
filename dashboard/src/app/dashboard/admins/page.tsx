'use client';

import { useState, useEffect } from 'react';
import axios from 'axios';

type Admin = {
  id: number;
  email: string;
  username: string;
  role: string;
  active: boolean;
};

export default function AdminsPage() {
  const [admins, setAdmins] = useState<Admin[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchAdmins = async () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) return;

        const response = await axios.get('http://localhost:8080/api/admin/users/admins', {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });

        setAdmins(response.data.data || []);
      } catch (err: any) {
        console.error('Error fetching admins:', err);
        setError(err.response?.data?.error || 'Không thể tải danh sách Admin');
      } finally {
        setIsLoading(false);
      }
    };

    fetchAdmins();
  }, []);

  const handleToggleStatus = async (adminId: number, currentStatus: boolean) => {
    try {
      const token = localStorage.getItem('token');
      if (!token) return;

      await axios.post(
        'http://localhost:8080/api/admin/users/update-status',
        {
          user_id: adminId,
          active: !currentStatus
        },
        {
          headers: {
            Authorization: `Bearer ${token}`
          }
        }
      );

      // Cập nhật trạng thái trong state
      setAdmins(admins.map(admin => 
        admin.id === adminId ? { ...admin, active: !currentStatus } : admin
      ));
    } catch (err: any) {
      console.error('Error updating status:', err);
      alert(err.response?.data?.error || 'Không thể cập nhật trạng thái');
    }
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
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Quản lý Admin</h1>
        <button className="btn btn-primary text-white">Thêm Admin</button>
      </div>
      
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <div className="overflow-x-auto">
            <table className="table table-zebra">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Tên người dùng</th>
                  <th>Email</th>
                  <th>Trạng thái</th>
                </tr>
              </thead>
              <tbody>
                {admins.length > 0 ? (
                  admins.map((admin) => (
                    <tr key={admin.id}>
                      <td>{admin.id}</td>
                      <td>{admin.username}</td>
                      <td>{admin.email}</td>
                      <td>
                        <button
                          className={`btn btn-sm ${admin.active ? 'btn-success' : 'btn-error'} text-white`}
                          onClick={() => handleToggleStatus(admin.id, admin.active)}
                        >
                          {admin.active ? 'Hoạt động' : 'Vô hiệu hóa'}
                        </button>
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={4} className="text-center py-4">
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
