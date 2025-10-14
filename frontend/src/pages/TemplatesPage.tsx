import { useState, useEffect } from 'react'
import { FileText, Plus, Trash2, Copy } from 'lucide-react'
import { api } from '@/lib/api'
import { toast } from 'sonner'

interface Template {
  id: string
  name: string
  description: string
  protocol: string
  config: Record<string, any>
  created_at: string
}

export default function TemplatesPage() {
  const [templates, setTemplates] = useState<Template[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchTemplates()
  }, [])

  const fetchTemplates = async () => {
    try {
      setLoading(true)
      const response = await api.get('/templates')
      setTemplates(response.data.data || [])
    } catch (error) {
      toast.error('Failed to fetch templates')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this template?')) return

    try {
      await api.delete(`/templates/${id}`)
      toast.success('Template deleted')
      fetchTemplates()
    } catch (error) {
      toast.error('Failed to delete template')
    }
  }

  const handleApply = async (id: string) => {
    try {
      await api.post(`/templates/${id}/apply`)
      toast.success('Template applied successfully')
    } catch (error) {
      toast.error('Failed to apply template')
    }
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
          <h1 className="text-3xl font-bold">Configuration Templates</h1>
          <p className="text-muted-foreground mt-1">
            Pre-configured templates for quick setup
          </p>
        </div>
        <button className="inline-flex items-center px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition">
          <Plus className="w-4 h-4 mr-2" />
          Create Template
        </button>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {templates.map((template) => (
          <div key={template.id} className="bg-card border border-border rounded-lg p-6">
            <div className="flex items-start justify-between mb-4">
              <div className="p-2 bg-primary/10 rounded">
                <FileText className="w-5 h-5 text-primary" />
              </div>
              <div className="flex gap-1">
                <button
                  onClick={() => handleApply(template.id)}
                  className="p-1.5 hover:bg-accent rounded"
                  title="Apply"
                >
                  <Copy className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleDelete(template.id)}
                  className="p-1.5 hover:bg-destructive/10 text-destructive rounded"
                  title="Delete"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
            
            <h3 className="font-semibold mb-1">{template.name}</h3>
            <p className="text-sm text-muted-foreground mb-3">{template.description}</p>
            
            <div className="flex items-center gap-2 text-xs">
              <span className="px-2 py-1 bg-muted rounded">{template.protocol}</span>
              <span className="text-muted-foreground">
                {new Date(template.created_at).toLocaleDateString()}
              </span>
            </div>
          </div>
        ))}

        {templates.length === 0 && (
          <div className="col-span-full text-center py-12">
            <FileText className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
            <p className="text-muted-foreground">No templates found</p>
          </div>
        )}
      </div>
    </div>
  )
}
