import { useState } from 'react'
import { User, Lock, Shield, Bell, Moon, Sun, Save } from 'lucide-react'
import { useAuthStore } from '@/stores/authStore'
import { api } from '@/lib/api'
import { toast } from 'sonner'

export default function SettingsPage() {
  const { user, setUser } = useAuthStore()
  const [activeTab, setActiveTab] = useState<'profile' | 'password' | '2fa' | 'preferences'>('profile')
  const [loading, setLoading] = useState(false)
  const [theme, setTheme] = useState<'light' | 'dark'>('light')

  // Profile form
  const [profileForm, setProfileForm] = useState({
    username: user?.username || '',
    email: user?.email || '',
  })

  // Password form
  const [passwordForm, setPasswordForm] = useState({
    current_password: '',
    new_password: '',
    confirm_password: '',
  })

  // 2FA state
  const [twoFAEnabled, setTwoFAEnabled] = useState(user?.two_factor_enabled || false)
  const [qrCode, setQrCode] = useState('')
  const [verifyCode, setVerifyCode] = useState('')
  const [showQR, setShowQR] = useState(false)

  const handleProfileUpdate = async (e: React.FormEvent) => {
    e.preventDefault()
    
    try {
      setLoading(true)
      await api.put(`/users/${user?.id}`, profileForm)
      setUser({ ...user, ...profileForm } as any)
      toast.success('Profile updated successfully')
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to update profile')
    } finally {
      setLoading(false)
    }
  }

  const handlePasswordChange = async (e: React.FormEvent) => {
    e.preventDefault()

    if (passwordForm.new_password !== passwordForm.confirm_password) {
      toast.error('Passwords do not match')
      return
    }

    if (passwordForm.new_password.length < 8) {
      toast.error('Password must be at least 8 characters')
      return
    }

    try {
      setLoading(true)
      await api.post('/auth/change-password', {
        current_password: passwordForm.current_password,
        new_password: passwordForm.new_password,
      })
      toast.success('Password changed successfully')
      setPasswordForm({ current_password: '', new_password: '', confirm_password: '' })
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to change password')
    } finally {
      setLoading(false)
    }
  }

  const handleSetup2FA = async () => {
    try {
      setLoading(true)
      const response = await api.post('/2fa/setup')
      setQrCode(response.data.qr_code)
      setShowQR(true)
      toast.info('Scan the QR code with your authenticator app')
    } catch (error) {
      toast.error('Failed to setup 2FA')
    } finally {
      setLoading(false)
    }
  }

  const handleEnable2FA = async () => {
    if (!verifyCode || verifyCode.length !== 6) {
      toast.error('Please enter a valid 6-digit code')
      return
    }

    try {
      setLoading(true)
      await api.post('/2fa/enable', { code: verifyCode })
      setTwoFAEnabled(true)
      setShowQR(false)
      setVerifyCode('')
      toast.success('2FA enabled successfully')
    } catch (error) {
      toast.error('Invalid verification code')
    } finally {
      setLoading(false)
    }
  }

  const handleDisable2FA = async () => {
    if (!confirm('Are you sure you want to disable 2FA?')) return

    try {
      setLoading(true)
      await api.post('/2fa/disable')
      setTwoFAEnabled(false)
      toast.success('2FA disabled successfully')
    } catch (error) {
      toast.error('Failed to disable 2FA')
    } finally {
      setLoading(false)
    }
  }

  const handleThemeChange = (newTheme: 'light' | 'dark') => {
    setTheme(newTheme)
    document.documentElement.classList.toggle('dark', newTheme === 'dark')
    localStorage.setItem('theme', newTheme)
    toast.success(`Theme changed to ${newTheme} mode`)
  }

  const tabs = [
    { id: 'profile', label: 'Profile', icon: User },
    { id: 'password', label: 'Password', icon: Lock },
    { id: '2fa', label: 'Two-Factor Auth', icon: Shield },
    { id: 'preferences', label: 'Preferences', icon: Bell },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-muted-foreground mt-1">
          Manage your account settings and preferences
        </p>
      </div>

      {/* Tabs */}
      <div className="border-b border-border">
        <div className="flex gap-2 overflow-x-auto">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as any)}
              className={`flex items-center gap-2 px-4 py-3 border-b-2 transition whitespace-nowrap ${
                activeTab === tab.id
                  ? 'border-primary text-primary'
                  : 'border-transparent text-muted-foreground hover:text-foreground'
              }`}
            >
              <tab.icon className="w-4 h-4" />
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Content */}
      <div className="bg-card border border-border rounded-lg p-6">
        {/* Profile Tab */}
        {activeTab === 'profile' && (
          <form onSubmit={handleProfileUpdate} className="space-y-6 max-w-2xl">
            <div>
              <h2 className="text-xl font-semibold mb-4">Profile Information</h2>
              <p className="text-sm text-muted-foreground mb-6">
                Update your account profile information
              </p>
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">Username</label>
              <input
                type="text"
                required
                value={profileForm.username}
                onChange={(e) => setProfileForm({ ...profileForm, username: e.target.value })}
                className="w-full px-4 py-2 border border-border rounded-md bg-background"
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">Email</label>
              <input
                type="email"
                required
                value={profileForm.email}
                onChange={(e) => setProfileForm({ ...profileForm, email: e.target.value })}
                className="w-full px-4 py-2 border border-border rounded-md bg-background"
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">Role</label>
              <input
                type="text"
                disabled
                value={user?.role || ''}
                className="w-full px-4 py-2 border border-border rounded-md bg-muted cursor-not-allowed"
              />
              <p className="text-xs text-muted-foreground mt-1">
                Contact an administrator to change your role
              </p>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="inline-flex items-center px-6 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition disabled:opacity-50"
            >
              <Save className="w-4 h-4 mr-2" />
              {loading ? 'Saving...' : 'Save Changes'}
            </button>
          </form>
        )}

        {/* Password Tab */}
        {activeTab === 'password' && (
          <form onSubmit={handlePasswordChange} className="space-y-6 max-w-2xl">
            <div>
              <h2 className="text-xl font-semibold mb-4">Change Password</h2>
              <p className="text-sm text-muted-foreground mb-6">
                Ensure your account is using a strong password
              </p>
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">Current Password</label>
              <input
                type="password"
                required
                value={passwordForm.current_password}
                onChange={(e) => setPasswordForm({ ...passwordForm, current_password: e.target.value })}
                className="w-full px-4 py-2 border border-border rounded-md bg-background"
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">New Password</label>
              <input
                type="password"
                required
                value={passwordForm.new_password}
                onChange={(e) => setPasswordForm({ ...passwordForm, new_password: e.target.value })}
                className="w-full px-4 py-2 border border-border rounded-md bg-background"
              />
              <p className="text-xs text-muted-foreground mt-1">
                Must be at least 8 characters
              </p>
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">Confirm New Password</label>
              <input
                type="password"
                required
                value={passwordForm.confirm_password}
                onChange={(e) => setPasswordForm({ ...passwordForm, confirm_password: e.target.value })}
                className="w-full px-4 py-2 border border-border rounded-md bg-background"
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="inline-flex items-center px-6 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition disabled:opacity-50"
            >
              <Lock className="w-4 h-4 mr-2" />
              {loading ? 'Changing...' : 'Change Password'}
            </button>
          </form>
        )}

        {/* 2FA Tab */}
        {activeTab === '2fa' && (
          <div className="space-y-6 max-w-2xl">
            <div>
              <h2 className="text-xl font-semibold mb-4">Two-Factor Authentication</h2>
              <p className="text-sm text-muted-foreground mb-6">
                Add an extra layer of security to your account
              </p>
            </div>

            <div className="p-4 border border-border rounded-lg">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium">Status</h3>
                  <p className="text-sm text-muted-foreground mt-1">
                    2FA is currently {twoFAEnabled ? 'enabled' : 'disabled'}
                  </p>
                </div>
                <div className={`px-3 py-1 rounded-full text-sm font-medium ${
                  twoFAEnabled
                    ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                    : 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200'
                }`}>
                  {twoFAEnabled ? 'Enabled' : 'Disabled'}
                </div>
              </div>
            </div>

            {!twoFAEnabled && !showQR && (
              <button
                onClick={handleSetup2FA}
                disabled={loading}
                className="inline-flex items-center px-6 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition disabled:opacity-50"
              >
                <Shield className="w-4 h-4 mr-2" />
                {loading ? 'Setting up...' : 'Setup 2FA'}
              </button>
            )}

            {showQR && (
              <div className="space-y-4">
                <div className="p-4 bg-white rounded-lg inline-block">
                  <img src={qrCode} alt="QR Code" className="w-48 h-48" />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">Verification Code</label>
                  <input
                    type="text"
                    maxLength={6}
                    value={verifyCode}
                    onChange={(e) => setVerifyCode(e.target.value.replace(/\D/g, ''))}
                    placeholder="Enter 6-digit code"
                    className="w-full px-4 py-2 border border-border rounded-md bg-background"
                  />
                </div>
                <button
                  onClick={handleEnable2FA}
                  disabled={loading}
                  className="inline-flex items-center px-6 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition disabled:opacity-50"
                >
                  {loading ? 'Verifying...' : 'Verify & Enable'}
                </button>
              </div>
            )}

            {twoFAEnabled && (
              <button
                onClick={handleDisable2FA}
                disabled={loading}
                className="inline-flex items-center px-6 py-2 bg-destructive text-destructive-foreground rounded-md hover:bg-destructive/90 transition disabled:opacity-50"
              >
                {loading ? 'Disabling...' : 'Disable 2FA'}
              </button>
            )}
          </div>
        )}

        {/* Preferences Tab */}
        {activeTab === 'preferences' && (
          <div className="space-y-6 max-w-2xl">
            <div>
              <h2 className="text-xl font-semibold mb-4">Preferences</h2>
              <p className="text-sm text-muted-foreground mb-6">
                Customize your experience
              </p>
            </div>

            <div className="p-4 border border-border rounded-lg">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium">Theme</h3>
                  <p className="text-sm text-muted-foreground mt-1">
                    Choose your preferred color scheme
                  </p>
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={() => handleThemeChange('light')}
                    className={`p-2 rounded-md border transition ${
                      theme === 'light'
                        ? 'border-primary bg-primary/10'
                        : 'border-border hover:bg-accent'
                    }`}
                  >
                    <Sun className="w-5 h-5" />
                  </button>
                  <button
                    onClick={() => handleThemeChange('dark')}
                    className={`p-2 rounded-md border transition ${
                      theme === 'dark'
                        ? 'border-primary bg-primary/10'
                        : 'border-border hover:bg-accent'
                    }`}
                  >
                    <Moon className="w-5 h-5" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
