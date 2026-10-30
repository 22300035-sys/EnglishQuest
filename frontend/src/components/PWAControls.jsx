import { useEffect, useState } from 'react'

export default function PWAControls() {
  const [prompt, setPrompt] = useState(null)
  useEffect(() => {
    const available = event => { event.preventDefault(); setPrompt(event) }
    const installed = () => setPrompt(null)
    window.addEventListener('beforeinstallprompt', available)
    window.addEventListener('appinstalled', installed)
    return () => { window.removeEventListener('beforeinstallprompt', available); window.removeEventListener('appinstalled', installed) }
  }, [])
  if (!prompt) return null
  return <button className="btn secondary small" onClick={async () => { await prompt.prompt(); await prompt.userChoice; setPrompt(null) }}>Instalar EnglishQuest</button>
}
