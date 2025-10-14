import { useState, useEffect } from 'react'
import { Plus, Download, Trash2, RefreshCw, Shield, Calendar } from 'lucide-react'
import { api } from '@/lib/api'
import { toast } from 'sonner'

interface Certificate {
  id: string
  domain: string
  type: 'acme' | 'self-signed' | 'custom'
  expires_at: string
  auto_renew: boolean
  status: 'active' | 'expired' | 'pending'
  created_at: string
}

export default function CertificatesPage() {
  const [certificates, setCertificates] = useState<Certificate[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [modalType, setModalType] = useState<'acme' | 'upload'>('acme')
  const [formData, setFormData] = useState({
    domain: '',
    email: '',
    auto_renew: true,
  })

  useEffect(() => {
    fetchCertificates()
  }, [])

  const fetchCertificates = async () => {
    try {
      setLoading(true)
      const response = await api.get('/certificates')
      setCertificates(response.data.data || [])
    } catch (error) {
      toast.error('Failed to fetch certificates')
    } finally {
      setLoading(false)
    }
  }

  const handleRequestACME = async (e: React.FormEvent) => {
    e.preventDefault()
    
    try {
      setLoading(true)
      await api.post('/certificates/acme', formData)
      toast.success('Certificate requested successfully')
      setShowModal(false)
      setFormData({ domain: '', email: '', auto_renew: true })
      fetchCertificates()
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to request certificate')
    } finally {
      setLoading(false)
    }
  }

  const handleRenew = async (id: string) => {
    try {
      await api.post(`/certificates/${id}/renew`)
      toast.success('Certificate renewed successfully')
      fetchCertificates()
    } catch (error) {
      toast.error('Failed to renew certificate')
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this certificate?')) return

    try {
      await api.delete(`/certificates/${id}`)
      toast.success('Certificate deleted successfully')
      fetchCertificates()
    } catch (error) {
      toast.error('Failed to delete certificate')
    }
  }

  const handleToggleAutoRenew = async (id: string, enabled: boolean) => {
    try {
      await api.patch(`/certificates/${id}`, { auto_renew: enabled })
      toast.success(`Auto-renew ${enabled ? 'enabled' : 'disabled'}`)
      fetchCertificates()
    } catch (error) {
      toast.error('Failed to update auto-renew')
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
      case 'expired': return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
      case 'pending': return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  const getDaysUntilExpiry = (expiresAt: string) => {
    const days = Math.ceil((new Date(expiresAt).getTime() - Date.now()) / (1000 * 60 * 60 * 24))
    return days > 0 ? `${days} days` : 'Expired'
  }

  if (loading && certificates.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">SSL Certificates</h1>
          <p className="text-muted-foreground mt-1">Manage SSL/TLS certificates</p>
        </div>
        <button
          onClick={() => { setModalType('acme'); setShowModal(true) }}
          className="inline-flex items-center px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition"
        >
          <Plus className="w-4 h-4 mr-2" />
          Request Certificate
        </button>
      </div>

      <div className="grid gap-4">
        {certificates.map((cert) => (
          <div key={cert.id} className="bg-card border border-border rounded-lg p-6">
            <div className="flex items-start justify-between mb-4">
              <div className="flex-1">
                <div className="flex items-center gap-3">
                  <Shield className="w-5 h-5 text-primary" />
                  <h3 className="text-lg font-semibold">{cert.domain}</h3>
                  <span className={`px-2 py-1 rounded text-xs font-medium ${getStatusColor(cert.status)}`}>
                    {cert.status.toUpperCase()}
                  </span>
                </div>
                <p className="text-sm text-muted-foreground mt-2">
                  Type: {cert.type} • Expires in {getDaysUntilExpiry(cert.expires_at)}
                </p>
              </div>
              <div className="flex gap-2">
                {cert.status === 'active' && (
                  <button
                    onClick={() => handleRenew(cert.id)}
                    className="p-2 hover:bg-accent rounded-md transition"
                    title="Renew"
                  >
                    <RefreshCw className="w-4 h-4" />
                  </button>
                )}
                <button
                  onClick={() => handleDelete(cert.id)}
                  className="p-2 hover:bg-destructive/10 text-destructive rounded-md transition"
                  title="Delete"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>

            <div className="flex items-center justify-between p-3 bg-muted rounded-md">
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={cert.auto_renew}
                  onChange={(e) => handleToggleAutoRenew(cert.id, e.target.checked)}
                  className="rounded"
                />
                <span className="text-sm">Auto-renew before expiration</span>
              </div>
              <span className="text-xs text-muted-foreground">
                <Calendar className="w-3 h-3 inline mr-1" />
                Created {new Date(cert.created_at).toLocaleDateString()}
              </span>
            </div>
          </div>
        ))}

        {certificates.length === 0 && (
          <div className="text-center py-12">
            <p className="text-muted-foreground">No certificates found</p>
          </div>
        )}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-card border border-border rounded-lg p-6 w-full max-w-md">
            <h2 className="text-xl font-bold mb-4">Request SSL Certificate</h2>
            <form onSubmit={handleRequestACME} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">Domain</label>
                <input
                  type="text"
                  required
                  value={formData.domain}
                  onChange={(e) => setFormData({ ...formData, domain: e.target.value })}
                  placeholder="example.com"
                  className="w-full px-3 py-2 border border-border rounded-md bg-background"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-2">Email</label>
                <input
                  type="email"
                  required
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  placeholder="admin@example.com"
                  className="w-full px-3 py-2 border border-border rounded-md bg-background"
                />
              </div>
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={formData.auto_renew}
                  onChange={(e) => setFormData({ ...formData, auto_renew: e.target.checked })}
                  className="rounded"
                />
                <label className="text-sm">Enable auto-renewal</label>
              </div>
              <div className="flex gap-2 justify-end">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="px-4 py-2 border border-border rounded-md hover:bg-accent transition"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={loading}
                  className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition disabled:opacity-50"
                >
                  {loading ? 'Requesting...' : 'Request'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
