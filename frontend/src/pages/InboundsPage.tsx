import { useState, useEffect } from 'react'
import { Plus, Edit, Trash2, Power, Link as LinkIcon, QrCode } from 'lucide-react'
import { api } from '../lib/api'
import { useAuthStore } from '../stores/authStore'

interface Inbound {
  id: string
  tag: string
  protocol: string
  port: number
  enable: boolean
  remark: string
  created_at: string
  settings: {
    clients: Array<{
      email: string
      id: string
      enable: boolean
    }>
  }
}

export default function InboundsPage() {
  const { user } = useAuthStore()
  const [inbounds, setInbounds] = useState<Inbound[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchInbounds()
  }, [])

  const fetchInbounds = async () => {
    try {
      setLoading(true)
      const response = await api.get('/inbounds')
      setInbounds(response.data.data)
    } catch (error) {
      console.error('Failed to fetch inbounds:', error)
    } finally {
      setLoading(false)
    }
  }

  const toggleInbound = async (id: string, enable: boolean) => {
    try {
      await api.post(`/inbounds/${id}/toggle`, { enable: !enable })
      fetchInbounds()
    } catch (error) {
      console.error('Failed to toggle inbound:', error)
    }
  }

  const deleteInbound = async (id: string) => {
    if (!confirm('Are you sure you want to delete this inbound?')) return

    try {
      await api.delete(`/inbounds/${id}`)
      fetchInbounds()
    } catch (error) {
      console.error('Failed to delete inbound:', error)
    }
  }

  const getProtocolColor = (protocol: string) => {
    const colors: Record<string, string> = {
      vmess: 'bg-blue-500',
      vless: 'bg-purple-500',
      trojan: 'bg-red-500',
      shadowsocks: 'bg-green-500',
      wireguard: 'bg-yellow-500',
    }
    return colors[protocol] || 'bg-gray-500'
  }

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
          <h1 className="text-3xl font-bold">Inbounds</h1>
          <p className="text-muted-foreground mt-1">
            Manage your Xray inbound configurations
          </p>
        </div>
        {user?.role === 'admin' && (
          <button
            onClick={() => setShowCreateModal(true)}
            className="inline-flex items-center px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition"
          >
            <Plus className="w-4 h-4 mr-2" />
            Create Inbound
          </button>
        )}
      </div>

      {/* Inbounds List */}
      <div className="grid gap-4">
        {inbounds.map((inbound) => (
          <div
            key={inbound.id}
            className="bg-card border border-border rounded-lg p-6 hover:shadow-md transition"
          >
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <div className="flex items-center gap-3">
                  <span
                    className={`${getProtocolColor(
                      inbound.protocol
                    )} text-white px-3 py-1 rounded-full text-sm font-medium uppercase`}
                  >
                    {inbound.protocol}
                  </span>
                  <h3 className="text-xl font-semibold">{inbound.tag}</h3>
                  <span
                    className={`px-2 py-1 rounded text-xs font-medium ${
                      inbound.enable
                        ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                        : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                    }`}
                  >
                    {inbound.enable ? 'Active' : 'Inactive'}
                  </span>
                </div>

                <div className="mt-3 grid grid-cols-3 gap-4 text-sm">
                  <div>
                    <span className="text-muted-foreground">Port:</span>
                    <span className="ml-2 font-medium">{inbound.port}</span>
                  </div>
                  <div>
                    <span className="text-muted-foreground">Clients:</span>
                    <span className="ml-2 font-medium">
                      {inbound.settings.clients?.length || 0}
                    </span>
                  </div>
                  <div>
                    <span className="text-muted-foreground">Created:</span>
                    <span className="ml-2 font-medium">
                      {new Date(inbound.created_at).toLocaleDateString()}
                    </span>
                  </div>
                </div>

                {inbound.remark && (
                  <p className="mt-2 text-sm text-muted-foreground">
                    {inbound.remark}
                  </p>
                )}
              </div>

              <div className="flex gap-2">
                <button
                  onClick={() => toggleInbound(inbound.id, inbound.enable)}
                  className="p-2 hover:bg-accent rounded-md transition"
                  title={inbound.enable ? 'Disable' : 'Enable'}
                >
                  <Power
                    className={`w-5 h-5 ${
                      inbound.enable ? 'text-green-600' : 'text-gray-400'
                    }`}
                  />
                </button>
                <button
                  className="p-2 hover:bg-accent rounded-md transition"
                  title="Get Links"
                >
                  <LinkIcon className="w-5 h-5" />
                </button>
                <button
                  className="p-2 hover:bg-accent rounded-md transition"
                  title="QR Code"
                >
                  <QrCode className="w-5 h-5" />
                </button>
                {user?.role === 'admin' && (
                  <>
                    <button
                      className="p-2 hover:bg-accent rounded-md transition"
                      title="Edit"
                    >
                      <Edit className="w-5 h-5" />
                    </button>
                    <button
                      onClick={() => deleteInbound(inbound.id)}
                      className="p-2 hover:bg-destructive/10 text-destructive rounded-md transition"
                      title="Delete"
                    >
                      <Trash2 className="w-5 h-5" />
                    </button>
                  </>
                )}
              </div>
            </div>

            {/* Clients Preview */}
            {inbound.settings.clients && inbound.settings.clients.length > 0 && (
              <div className="mt-4 pt-4 border-t border-border">
                <h4 className="text-sm font-medium mb-2">Clients:</h4>
                <div className="flex flex-wrap gap-2">
                  {inbound.settings.clients.slice(0, 5).map((client) => (
                    <span
                      key={client.id}
                      className="px-2 py-1 bg-accent rounded text-xs"
                    >
                      {client.email}
                    </span>
                  ))}
                  {inbound.settings.clients.length > 5 && (
                    <span className="px-2 py-1 bg-accent rounded text-xs">
                      +{inbound.settings.clients.length - 5} more
                    </span>
                  )}
                </div>
              </div>
            )}
          </div>
        ))}

        {inbounds.length === 0 && (
          <div className="text-center py-12">
            <p className="text-muted-foreground">No inbounds configured yet</p>
            {user?.role === 'admin' && (
              <button
                onClick={() => setShowCreateModal(true)}
                className="mt-4 inline-flex items-center px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition"
              >
                <Plus className="w-4 h-4 mr-2" />
                Create Your First Inbound
              </button>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
