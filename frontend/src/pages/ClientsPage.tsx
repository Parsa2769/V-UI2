import { useState, useEffect } from 'react'
import { Plus, Trash2, Link as LinkIcon, QrCode, Download, RefreshCw, Search, Filter } from 'lucide-react'
import { api } from '@/lib/api'
import { toast } from 'sonner'

interface Client {
  id: string
  email: string
  total_gb: number
  upload: number
  download: number
  expiry_time: number
  limit_ip: number
  enable: boolean
  created_at: string
  subscription_url?: string
}

export default function ClientsPage() {
  const [clients, setClients] = useState<Client[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [showLinkModal, setShowLinkModal] = useState(false)
  const [selectedClient, setSelectedClient] = useState<Client | null>(null)
  const [shareLink, setShareLink] = useState('')
  const [searchTerm, setSearchTerm] = useState('')
  const [filterStatus, setFilterStatus] = useState<'all' | 'active' | 'expired'>('all')
  const [formData, setFormData] = useState({
    email: '',
    total_gb: 50,
    limit_ip: 2,
    expiry_days: 30,
  })

  useEffect(() => {
    fetchClients()
  }, [])

  const fetchClients = async () => {
    try {
      setLoading(true)
      const response = await api.get('/clients')
      setClients(response.data.data || [])
    } catch (error) {
      toast.error('Failed to fetch clients')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    try {
      if (selectedClient) {
        await api.put(`/clients/${selectedClient.id}`, formData)
        toast.success('Client updated successfully')
      } else {
        await api.post('/clients', {
          ...formData,
          expiry_time: Date.now() + (formData.expiry_days * 24 * 60 * 60 * 1000)
        })
        toast.success('Client created successfully')
      }
      
      setShowModal(false)
      setSelectedClient(null)
      setFormData({ email: '', total_gb: 50, limit_ip: 2, expiry_days: 30 })
      fetchClients()
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Operation failed')
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this client?')) return

    try {
      await api.delete(`/clients/${id}`)
      toast.success('Client deleted successfully')
      fetchClients()
    } catch (error) {
      toast.error('Failed to delete client')
    }
  }

  const handleResetTraffic = async (id: string) => {
    if (!confirm('Reset traffic for this client?')) return

    try {
      await api.post(`/clients/${id}/reset-traffic`)
      toast.success('Traffic reset successfully')
      fetchClients()
    } catch (error) {
      toast.error('Failed to reset traffic')
    }
  }

  const handleGetLink = async (client: Client) => {
    try {
      const response = await api.get(`/clients/${client.id}/link`)
      setShareLink(response.data.link || 'No link available')
      setSelectedClient(client)
      setShowLinkModal(true)
    } catch (error) {
      toast.error('Failed to get share link')
    }
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    toast.success('Copied to clipboard!')
  }

  const calculateUsage = (client: Client) => {
    const used = (client.upload + client.download) / (1024 * 1024 * 1024) // Convert to GB
    const total = client.total_gb
    const percentage = total > 0 ? (used / total) * 100 : 0
    return { used: used.toFixed(2), total, percentage: Math.min(percentage, 100) }
  }

  const isExpired = (expiryTime: number) => {
    return expiryTime > 0 && expiryTime < Date.now()
  }

  const getDaysRemaining = (expiryTime: number) => {
    if (expiryTime === 0) return 'Unlimited'
    const days = Math.ceil((expiryTime - Date.now()) / (1000 * 60 * 60 * 24))
    return days > 0 ? `${days} days` : 'Expired'
  }

  const filteredClients = clients.filter(client => {
    const matchesSearch = client.email.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesFilter = 
      filterStatus === 'all' ||
      (filterStatus === 'active' && client.enable && !isExpired(client.expiry_time)) ||
      (filterStatus === 'expired' && (isExpired(client.expiry_time) || !client.enable))
    return matchesSearch && matchesFilter
  })

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Clients Management</h1>
          <p className="text-muted-foreground mt-1">
            Manage your VPN clients and monitor their usage
          </p>
        </div>
        <button
          onClick={() => {
            setSelectedClient(null)
            setFormData({ email: '', total_gb: 50, limit_ip: 2, expiry_days: 30 })
            setShowModal(true)
          }}
          className="inline-flex items-center px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition"
        >
          <Plus className="w-4 h-4 mr-2" />
          Add Client
        </button>
      </div>

      {/* Search and Filter */}
      <div className="flex gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
          <input
            type="text"
            placeholder="Search clients..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full pl-10 pr-4 py-2 border border-border rounded-md bg-background"
          />
        </div>
        <div className="flex items-center gap-2">
          <Filter className="w-4 h-4 text-muted-foreground" />
          <select
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value as any)}
            className="px-3 py-2 border border-border rounded-md bg-background"
          >
            <option value="all">All</option>
            <option value="active">Active</option>
            <option value="expired">Expired</option>
          </select>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-card border border-border rounded-lg p-4">
          <div className="text-sm text-muted-foreground">Total Clients</div>
          <div className="text-2xl font-bold mt-1">{clients.length}</div>
        </div>
        <div className="bg-card border border-border rounded-lg p-4">
          <div className="text-sm text-muted-foreground">Active Clients</div>
          <div className="text-2xl font-bold mt-1 text-green-600">
            {clients.filter(c => c.enable && !isExpired(c.expiry_time)).length}
          </div>
        </div>
        <div className="bg-card border border-border rounded-lg p-4">
          <div className="text-sm text-muted-foreground">Expired Clients</div>
          <div className="text-2xl font-bold mt-1 text-red-600">
            {clients.filter(c => isExpired(c.expiry_time) || !c.enable).length}
          </div>
        </div>
      </div>

      {/* Clients Grid */}
      <div className="grid gap-4">
        {filteredClients.map((client) => {
          const usage = calculateUsage(client)
          const expired = isExpired(client.expiry_time)
          
          return (
            <div
              key={client.id}
              className="bg-card border border-border rounded-lg p-6 hover:shadow-md transition"
            >
              <div className="flex items-start justify-between mb-4">
                <div className="flex-1">
                  <div className="flex items-center gap-3">
                    <h3 className="text-lg font-semibold">{client.email}</h3>
                    <span
                      className={`px-2 py-1 rounded text-xs font-medium ${
                        client.enable && !expired
                          ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                          : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                      }`}
                    >
                      {client.enable && !expired ? 'Active' : 'Inactive'}
                    </span>
                  </div>
                  <div className="mt-2 text-sm text-muted-foreground">
                    Created: {new Date(client.created_at).toLocaleDateString()}
                  </div>
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={() => handleGetLink(client)}
                    className="p-2 hover:bg-accent rounded-md transition"
                    title="Get Link"
                  >
                    <LinkIcon className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => handleResetTraffic(client.id)}
                    className="p-2 hover:bg-accent rounded-md transition"
                    title="Reset Traffic"
                  >
                    <RefreshCw className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => handleDelete(client.id)}
                    className="p-2 hover:bg-destructive/10 text-destructive rounded-md transition"
                    title="Delete"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>

              {/* Traffic Usage */}
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Traffic Usage</span>
                  <span className="font-medium">
                    {usage.used} GB / {usage.total} GB
                  </span>
                </div>
                <div className="w-full bg-muted rounded-full h-2">
                  <div
                    className={`h-2 rounded-full transition-all ${
                      usage.percentage > 90
                        ? 'bg-red-500'
                        : usage.percentage > 70
                        ? 'bg-yellow-500'
                        : 'bg-green-500'
                    }`}
                    style={{ width: `${usage.percentage}%` }}
                  />
                </div>
              </div>

              {/* Details Grid */}
              <div className="grid grid-cols-3 gap-4 mt-4 text-sm">
                <div>
                  <div className="text-muted-foreground">Expiry</div>
                  <div className="font-medium">{getDaysRemaining(client.expiry_time)}</div>
                </div>
                <div>
                  <div className="text-muted-foreground">IP Limit</div>
                  <div className="font-medium">{client.limit_ip || 'Unlimited'}</div>
                </div>
                <div>
                  <div className="text-muted-foreground">Download</div>
                  <div className="font-medium">
                    {(client.download / (1024 * 1024 * 1024)).toFixed(2)} GB
                  </div>
                </div>
              </div>
            </div>
          )
        })}

        {filteredClients.length === 0 && (
          <div className="text-center py-12">
            <p className="text-muted-foreground">No clients found</p>
          </div>
        )}
      </div>

      {/* Create/Edit Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-card border border-border rounded-lg p-6 w-full max-w-md">
            <h2 className="text-xl font-bold mb-4">
              {selectedClient ? 'Edit Client' : 'Add Client'}
            </h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">Email</label>
                <input
                  type="email"
                  required
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  className="w-full px-3 py-2 border border-border rounded-md bg-background"
                  placeholder="client@example.com"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Traffic Limit (GB)</label>
                <input
                  type="number"
                  required
                  min="1"
                  value={formData.total_gb}
                  onChange={(e) => setFormData({ ...formData, total_gb: parseInt(e.target.value) })}
                  className="w-full px-3 py-2 border border-border rounded-md bg-background"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">IP Limit</label>
                <input
                  type="number"
                  required
                  min="0"
                  value={formData.limit_ip}
                  onChange={(e) => setFormData({ ...formData, limit_ip: parseInt(e.target.value) })}
                  className="w-full px-3 py-2 border border-border rounded-md bg-background"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Expiry (Days)</label>
                <input
                  type="number"
                  required
                  min="1"
                  value={formData.expiry_days}
                  onChange={(e) => setFormData({ ...formData, expiry_days: parseInt(e.target.value) })}
                  className="w-full px-3 py-2 border border-border rounded-md bg-background"
                />
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
                  className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition"
                >
                  {selectedClient ? 'Update' : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Link Modal */}
      {showLinkModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-card border border-border rounded-lg p-6 w-full max-w-lg">
            <h2 className="text-xl font-bold mb-4">Share Link</h2>
            <div className="space-y-4">
              <div className="p-4 bg-muted rounded-md break-all font-mono text-sm">
                {shareLink}
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => copyToClipboard(shareLink)}
                  className="flex-1 inline-flex items-center justify-center px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition"
                >
                  <Download className="w-4 h-4 mr-2" />
                  Copy Link
                </button>
                <button
                  onClick={() => {/* QR Code generation */}}
                  className="flex-1 inline-flex items-center justify-center px-4 py-2 border border-border rounded-md hover:bg-accent transition"
                >
                  <QrCode className="w-4 h-4 mr-2" />
                  Show QR
                </button>
              </div>
              <button
                onClick={() => setShowLinkModal(false)}
                className="w-full px-4 py-2 border border-border rounded-md hover:bg-accent transition"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
