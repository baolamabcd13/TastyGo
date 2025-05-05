'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Sidebar from '@/components/layout/Sidebar';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();

  useEffect(() => {
    // Kiểm tra xem người dùng đã đăng nhập chưa
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/login');
    }
  }, [router]);

  return (
    <div className="flex h-screen bg-background">
      <Sidebar />
      
      <div className="flex-1 lg:ml-64">
        <main className="p-4 lg:p-8 min-h-screen">
          {children}
        </main>
      </div>
    </div>
  );
}
