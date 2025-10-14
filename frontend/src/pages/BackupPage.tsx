import { useState, useEffect } from 'react'
import { Download, Upload, Trash2, Database, HardDrive } from 'lucide-react'
import { api } from '@/lib/api'
import { toast } from 'sonner'

interface Backup {
  filename: string
  size: number
  created_at: string
}

export default function BackupPage() {
  const [backups, setBackups] = useState<Backup[]>([])
  const [loading, setLoading] = useState(true)
  const [creating, setCreating] = useState(false)

  useEffect(() => {
    fetchBackups()
  }, [])

  const fetchBackups = async () => {
    try {
      setLoading(true)
      const response = await api.get('/backup/list')
      setBackups(response.data || [])
    } catch (error) {
      toast.error('Failed to fetch backups')
    } finally {
      setLoading(false)
    }
  }

  const handleCreateBackup = async () => {
    try {
      setCreating(true)
      const response = await api.post('/backup/create')
      toast.success(`Backup created: ${response.data.filename}`)
      fetchBackups()
    } catch (error) {
      toast.error('Failed to create backup')
    } finally {
      setCreating(false)
    }
  }

  const handleDownload = async (filename: string) => {
    try {
      const response = await api.get(`/backup/download/${filename}`, {
        responseType: 'blob',
      })
      
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', filename)
      document.body.appendChild(link)
      link.click()
      link.remove()
      
      toast.success('Backup downloaded')
    } catch (error) {
      toast.error('Failed to download backup')
    }
  }

  const handleRestore = async (filename: string) => {
    if (!confirm(`Are you sure you want to restore from ${filename}? This will overwrite current data.`)) return

    try {
      await api.post('/backup/restore', { filename })
      toast.success('Backup restored successfully')
    } catch (error) {
      toast.error('Failed to restore backup')
    }
  }

  const handleDelete = async (filename: string) => {
    if (!confirm(`Delete backup ${filename}?`)) return

    try {
      await api.delete(`/backup/${filename}`)
      toast.success('Backup deleted')
      fetchBackups()
    } catch (error) {
      toast.error('Failed to delete backup')
    }
  }

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
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
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Backup & Restore</h1>
          <p className="text-muted-foreground mt-1">
            Create and manage database backups
          </p>
        </div>
        <button
          onClick={handleCreateBackup}
          disabled={creating}
          className="inline-flex items-center px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition disabled:opacity-50"
        >
          <Database className="w-4 h-4 mr-2" />
          {creating ? 'Creating...' : 'Create Backup'}
        </button>
      </div>

      <div className="grid gap-4">
        {backups.map((backup) => (
          <div
            key={backup.filename}
            className="bg-card border border-border rounded-lg p-6 hover:shadow-md transition"
          >
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-4 flex-1">
                <div className="p-3 bg-primary/10 rounded-lg">
                  <HardDrive className="w-6 h-6 text-primary" />
                </div>
                <div className="flex-1 min-w-0">
                  <h3 className="font-semibold truncate">{backup.filename}</h3>
                  <div className="flex items-center gap-4 mt-1 text-sm text-muted-foreground">
                    <span>{formatBytes(backup.size)}</span>
                    <span>•</span>
                    <span>{new Date(backup.created_at).toLocaleString()}</span>
                  </div>
                </div>
              </div>
              
              <div className="flex gap-2">
                <button
                  onClick={() => handleDownload(backup.filename)}
                  className="p-2 hover:bg-accent rounded-md transition"
                  title="Download"
                >
                  <Download className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleRestore(backup.filename)}
                  className="p-2 hover:bg-accent rounded-md transition"
                  title="Restore"
                >
                  <Upload className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleDelete(backup.filename)}
                  className="p-2 hover:bg-destructive/10 text-destructive rounded-md transition"
                  title="Delete"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        ))}

        {backups.length === 0 && (
          <div className="text-center py-12">
            <Database className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
            <p className="text-muted-foreground">No backups found</p>
            <p className="text-sm text-muted-foreground mt-1">Create your first backup to get started</p>
          </div>
        )}
      </div>

      <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
        <h3 className="font-semibold text-yellow-900 dark:text-yellow-200 mb-2">⚠️ Important Notes</h3>
        <ul className="text-sm text-yellow-800 dark:text-yellow-300 space-y-1">
          <li>• Backups include all users, clients, inbounds, and configurations</li>
          <li>• Restoring a backup will overwrite all current data</li>
          <li>• Download backups regularly and store them securely</li>
          <li>• Backups do not include SSL certificates or system files</li>
        </ul>
      </div>
    </div>
  )
}
