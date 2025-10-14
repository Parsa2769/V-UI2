import { useState, useEffect } from 'react'
import { Bot, Save, TestTube, Users, Bell } from 'lucide-react'
import { api } from '@/lib/api'
import { toast } from 'sonner'

export default function TelegramPage() {
  const [loading, setLoading] = useState(false)
  const [config, setConfig] = useState({
    bot_token: '',
    admin_ids: '',
    enabled: false,
    daily_report: true,
    traffic_alerts: true,
  })

  useEffect(() => {
    fetchConfig()
  }, [])

  const fetchConfig = async () => {
    try {
      const response = await api.get('/telegram/config')
      const data = response.data
      setConfig({
        bot_token: data.bot_token || '',
        admin_ids: data.admin_ids?.join(', ') || '',
        enabled: data.enabled || false,
        daily_report: data.daily_report || true,
        traffic_alerts: data.traffic_alerts || true,
      })
    } catch (error) {
      console.error('Failed to fetch config')
    }
  }

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault()
    
    try {
      setLoading(true)
      await api.post('/telegram/config', {
        bot_token: config.bot_token,
        admin_ids: config.admin_ids.split(',').map(id => id.trim()).filter(Boolean),
        enabled: config.enabled,
        daily_report: config.daily_report,
        traffic_alerts: config.traffic_alerts,
      })
      toast.success('Telegram settings saved successfully')
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to save settings')
    } finally {
      setLoading(false)
    }
  }

  const handleTest = async () => {
    try {
      setLoading(true)
      await api.post('/telegram/test')
      toast.success('Test message sent! Check your Telegram')
    } catch (error) {
      toast.error('Failed to send test message')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Telegram Bot</h1>
        <p className="text-muted-foreground mt-1">
          Configure Telegram bot for notifications and management
        </p>
      </div>

      <form onSubmit={handleSave} className="space-y-6 max-w-2xl">
        <div className="bg-card border border-border rounded-lg p-6 space-y-6">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-lg font-semibold">Bot Status</h2>
              <p className="text-sm text-muted-foreground">Enable or disable the Telegram bot</p>
            </div>
            <label className="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                checked={config.enabled}
                onChange={(e) => setConfig({ ...config, enabled: e.target.checked })}
                className="sr-only peer"
              />
              <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-primary"></div>
            </label>
          </div>

          <div>
            <label className="block text-sm font-medium mb-2 flex items-center gap-2">
              <Bot className="w-4 h-4" />
              Bot Token
            </label>
            <input
              type="password"
              required
              value={config.bot_token}
              onChange={(e) => setConfig({ ...config, bot_token: e.target.value })}
              placeholder="123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
              className="w-full px-4 py-2 border border-border rounded-md bg-background font-mono text-sm"
            />
            <p className="text-xs text-muted-foreground mt-1">
              Get your bot token from @BotFather on Telegram
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium mb-2 flex items-center gap-2">
              <Users className="w-4 h-4" />
              Admin User IDs
            </label>
            <input
              type="text"
              required
              value={config.admin_ids}
              onChange={(e) => setConfig({ ...config, admin_ids: e.target.value })}
              placeholder="123456789, 987654321"
              className="w-full px-4 py-2 border border-border rounded-md bg-background"
            />
            <p className="text-xs text-muted-foreground mt-1">
              Comma-separated list of Telegram user IDs who can control the bot
            </p>
          </div>

          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <Bell className="w-4 h-4" />
              Notifications
            </h3>
            
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={config.daily_report}
                onChange={(e) => setConfig({ ...config, daily_report: e.target.checked })}
                className="rounded"
              />
              <label className="text-sm">Send daily traffic report</label>
            </div>

            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={config.traffic_alerts}
                onChange={(e) => setConfig({ ...config, traffic_alerts: e.target.checked })}
                className="rounded"
              />
              <label className="text-sm">Send traffic limit alerts</label>
            </div>
          </div>
        </div>

        <div className="flex gap-3">
          <button
            type="submit"
            disabled={loading}
            className="inline-flex items-center px-6 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition disabled:opacity-50"
          >
            <Save className="w-4 h-4 mr-2" />
            {loading ? 'Saving...' : 'Save Settings'}
          </button>
          
          <button
            type="button"
            onClick={handleTest}
            disabled={loading || !config.enabled}
            className="inline-flex items-center px-6 py-2 border border-border rounded-md hover:bg-accent transition disabled:opacity-50"
          >
            <TestTube className="w-4 h-4 mr-2" />
            Send Test Message
          </button>
        </div>
      </form>

      <div className="bg-card border border-border rounded-lg p-6 max-w-2xl">
        <h3 className="font-semibold mb-3">Available Commands</h3>
        <div className="space-y-2 text-sm font-mono">
          <div><span className="text-primary">/start</span> - Start the bot</div>
          <div><span className="text-primary">/status</span> - Get system status</div>
          <div><span className="text-primary">/traffic</span> - Get traffic statistics</div>
          <div><span className="text-primary">/users</span> - List all users</div>
          <div><span className="text-primary">/inbounds</span> - List all inbounds</div>
          <div><span className="text-primary">/restart</span> - Restart Xray service</div>
          <div><span className="text-primary">/help</span> - Show help message</div>
        </div>
      </div>
    </div>
  )
}
