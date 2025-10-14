import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Toaster } from 'sonner'
import { useAuthStore } from '@/stores/authStore'

// Pages
import LoginPage from '@/pages/LoginPage'
import DashboardPage from '@/pages/DashboardPage'
import UsersPage from '@/pages/UsersPage'
import ClientsPage from '@/pages/ClientsPage'
import InboundsPage from '@/pages/InboundsPage'
import NodesPage from '@/pages/NodesPage'
import TrafficPage from '@/pages/TrafficPage'
import AuditLogsPage from '@/pages/AuditLogsPage'
import SettingsPage from '@/pages/SettingsPage'
import CertificatesPage from '@/pages/CertificatesPage'
import TelegramPage from '@/pages/TelegramPage'
import BackupPage from '@/pages/BackupPage'
import TemplatesPage from '@/pages/TemplatesPage'

// Layout
import MainLayout from '@/components/layout/MainLayout'

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />
}

function App() {
  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          
          <Route
            path="/"
            element={
              <PrivateRoute>
                <MainLayout />
              </PrivateRoute>
            }
          >
            <Route index element={<DashboardPage />} />
            <Route path="users" element={<UsersPage />} />
            <Route path="clients" element={<ClientsPage />} />
            <Route path="inbounds" element={<InboundsPage />} />
            <Route path="nodes" element={<NodesPage />} />
            <Route path="traffic" element={<TrafficPage />} />
            <Route path="certificates" element={<CertificatesPage />} />
            <Route path="telegram" element={<TelegramPage />} />
            <Route path="backup" element={<BackupPage />} />
            <Route path="templates" element={<TemplatesPage />} />
            <Route path="audit" element={<AuditLogsPage />} />
            <Route path="settings" element={<SettingsPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
      
      <Toaster position="top-right" richColors />
    </>
  )
}

export default App
