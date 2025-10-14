import { useState, useEffect } from 'react'
import { Plus, Edit, Trash2, Server, Activity, AlertCircle, CheckCircle, Search } from 'lucide-react'
import { api } from '@/lib/api'
import { toast } from 'sonner'

interface Node {
  id: string
  name: string
  address: string
  port: number
  api_port: number
  status: 'online' | 'offline' | 'error'
  load: number
  memory_usage: number
  cpu_usage: number
  uptime: number
  created_at: string
}

export default function NodesPage() {
  const [nodes, setNodes] = useState<Node[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [selectedNode, setSelectedNode] = useState<Node | null>(null)
  const [searchTerm, setSearchTerm] = useState('')
  const [formData, setFormData] = useState({
    name: '',
    address: '',
    port: 443,
    api_port: 10085,
  })

  useEffect(() => {
    fetchNodes()
    const interval = setInterval(fetchNodes, 30000) // Update every 30s
    return () => clearInterval(interval)
  }, [])

  const fetchNodes = async () => {
    try {
      setLoading(true)
      const response = await api.get('/nodes')
      setNodes(response.data.data || [])
    } catch (error) {
      toast.error('Failed to fetch nodes')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    try {
      if (selectedNode) {
        await api.put(`/nodes/${selectedNode.id}`, formData)
        toast.success('Node updated successfully')
      } else {
        await api.post('/nodes', formData)
        toast.success('Node created successfully')
      }
      
      setShowModal(false)
      setSelectedNode(null)
      setFormData({ name: '', address: '', port: 443, api_port: 10085 })
      fetchNodes()
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Operation failed')
    }
  }

  const handleEdit = (node: Node) => {
    setSelectedNode(node)
    setFormData({
      name: node.name,
      address: node.address,
      port: node.port,
      api_port: node.api_port,
    })
    setShowModal(true)
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this node?')) return

    try {
      await api.delete(`/nodes/${id}`)
      toast.success('Node deleted successfully')
      fetchNodes()
    } catch (error) {
      toast.error('Failed to delete node')
    }
  }

  const handleTestConnection = async (id: string) => {
    try {
      toast.info('Testing connection...')
      await api.post(`/nodes/${id}/test`)
      toast.success('Connection successful!')
      fetchNodes()
    } catch (error) {
      toast.error('Connection failed')
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'online':
        return 'text-green-600 bg-green-100 dark:bg-green-900 dark:text-green-200'
      case 'offline':
        return 'text-gray-600 bg-gray-100 dark:bg-gray-900 dark:text-gray-200'
      case 'error':
        return 'text-red-600 bg-red-100 dark:bg-red-900 dark:text-red-200'
      default:
        return 'text-gray-600 bg-gray-100'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'online':
        return <CheckCircle className="w-5 h-5" />
      case 'offline':
        return <AlertCircle className="w-5 h-5" />
      case 'error':
        return <AlertCircle className="w-5 h-5" />
      default:
        return <Server className="w-5 h-5" />
    }
  }

  const formatUptime = (seconds: number) => {
    const days = Math.floor(seconds / 86400)
    const hours = Math.floor((seconds % 86400) / 3600)
    if (days > 0) return `${days}d ${hours}h`
    return `${hours}h`
  }

  const filteredNodes = nodes.filter(node =>
    node.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    node.address.toLowerCase().includes(searchTerm.toLowerCase())
  )

  if (loading && nodes.length === 0) {
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
          <h1 className="text-3xl font-bold">Nodes Management</h1>
          <p className="text-muted-foreground mt-1">
            Manage your Xray server nodes
          </p>
        </div>
        <button
          onClick={() => {
            setSelectedNode(null)
            setFormData({ name: '', address: '', port: 443, api_port: 10085 })
            setShowModal(true)
          }}
          className="inline-flex items-center px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition"
        >
          <Plus className="w-4 h-4 mr-2" />
          Add Node
        </button>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
        <input
          type="text"
          placeholder="Search nodes..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full pl-10 pr-4 py-2 border border-border rounded-md bg-background"
        />
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-card border border-border rounded-lg p-4">
          <div className="text-sm text-muted-foreground">Total Nodes</div>
          <div className="text-2xl font-bold mt-1">{nodes.length}</div>
        </div>
        <div className="bg-card border border-border rounded-lg p-4">
          <div className="text-sm text-muted-foreground">Online</div>
          <div className="text-2xl font-bold mt-1 text-green-600">
            {nodes.filter(n => n.status === 'online').length}
          </div>
        </div>
        <div className="bg-card border border-border rounded-lg p-4">
          <div className="text-sm text-muted-foreground">Offline</div>
          <div className="text-2xl font-bold mt-1 text-gray-600">
            {nodes.filter(n => n.status === 'offline').length}
          </div>
        </div>
        <div className="bg-card border border-border rounded-lg p-4">
          <div className="text-sm text-muted-foreground">Errors</div>
          <div className="text-2xl font-bold mt-1 text-red-600">
            {nodes.filter(n => n.status === 'error').length}
          </div>
        </div>
      </div>

      {/* Nodes Grid */}
      <div className="grid gap-4">
        {filteredNodes.map((node) => (
          <div
            key={node.id}
            className="bg-card border border-border rounded-lg p-6 hover:shadow-md transition"
          >
            <div className="flex items-start justify-between mb-4">
              <div className="flex-1">
                <div className="flex items-center gap-3">
                  <Server className="w-5 h-5 text-muted-foreground" />
                  <h3 className="text-lg font-semibold">{node.name}</h3>
                  <span className={`flex items-center gap-1 px-2 py-1 rounded text-xs font-medium ${getStatusColor(node.status)}`}>
                    {getStatusIcon(node.status)}
                    {node.status.toUpperCase()}
                  </span>
                </div>
                <div className="mt-2 text-sm text-muted-foreground">
                  {node.address}:{node.port}
                </div>
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => handleTestConnection(node.id)}
                  className="p-2 hover:bg-accent rounded-md transition"
                  title="Test Connection"
                >
                  <Activity className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleEdit(node)}
                  className="p-2 hover:bg-accent rounded-md transition"
                  title="Edit"
                >
                  <Edit className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleDelete(node.id)}
                  className="p-2 hover:bg-destructive/10 text-destructive rounded-md transition"
                  title="Delete"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>

            {/* Metrics */}
            {node.status === 'online' && (
              <div className="grid grid-cols-4 gap-4 mt-4">
                <div>
                  <div className="text-xs text-muted-foreground mb-1">CPU Usage</div>
                  <div className="flex items-center gap-2">
                    <div className="flex-1 bg-muted rounded-full h-2">
                      <div
                        className={`h-2 rounded-full ${node.cpu_usage > 80 ? 'bg-red-500' : node.cpu_usage > 60 ? 'bg-yellow-500' : 'bg-green-500'}`}
                        style={{ width: `${node.cpu_usage}%` }}
                      />
                    </div>
                    <span className="text-xs font-medium">{node.cpu_usage}%</span>
                  </div>
                </div>
                <div>
                  <div className="text-xs text-muted-foreground mb-1">Memory</div>
                  <div className="flex items-center gap-2">
                    <div className="flex-1 bg-muted rounded-full h-2">
                      <div
                        className={`h-2 rounded-full ${node.memory_usage > 80 ? 'bg-red-500' : node.memory_usage > 60 ? 'bg-yellow-500' : 'bg-green-500'}`}
                        style={{ width: `${node.memory_usage}%` }}
                      />
                    </div>
                    <span className="text-xs font-medium">{node.memory_usage}%</span>
                  </div>
                </div>
                <div>
                  <div className="text-xs text-muted-foreground mb-1">Load</div>
                  <div className="text-sm font-medium">{node.load.toFixed(2)}</div>
                </div>
                <div>
                  <div className="text-xs text-muted-foreground mb-1">Uptime</div>
                  <div className="text-sm font-medium">{formatUptime(node.uptime)}</div>
                </div>
              </div>
            )}
          </div>
        ))}

        {filteredNodes.length === 0 && (
          <div className="text-center py-12">
            <p className="text-muted-foreground">No nodes found</p>
          </div>
        )}
      </div>

      {/* Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-card border border-border rounded-lg p-6 w-full max-w-md">
            <h2 className="text-xl font-bold mb-4">
              {selectedNode ? 'Edit Node' : 'Add Node'}
            </h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">Name</label>
                <input
                  type="text"
                  required
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="w-full px-3 py-2 border border-border rounded-md bg-background"
                  placeholder="My Node"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Address</label>
                <input
                  type="text"
                  required
                  value={formData.address}
                  onChange={(e) => setFormData({ ...formData, address: e.target.value })}
                  className="w-full px-3 py-2 border border-border rounded-md bg-background"
                  placeholder="example.com or 1.2.3.4"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-1">Port</label>
                  <input
                    type="number"
                    required
                    min="1"
                    max="65535"
                    value={formData.port}
                    onChange={(e) => setFormData({ ...formData, port: parseInt(e.target.value) })}
                    className="w-full px-3 py-2 border border-border rounded-md bg-background"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">API Port</label>
                  <input
                    type="number"
                    required
                    min="1"
                    max="65535"
                    value={formData.api_port}
                    onChange={(e) => setFormData({ ...formData, api_port: parseInt(e.target.value) })}
                    className="w-full px-3 py-2 border border-border rounded-md bg-background"
                  />
                </div>
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
                  {selectedNode ? 'Update' : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
