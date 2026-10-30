import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.jsx'
import { BrowserRouter } from 'react-router-dom'
import AuthProvider from './auth/AuthProvider.jsx'
import './quest.css'

// Service workers are deliberately disabled in development to avoid stale HMR.
if (import.meta.env.PROD && 'serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js', { updateViaCache: 'none' }).catch(() => {
      console.warn('La instalación offline no está disponible. La aplicación online puede seguir utilizándose.')
    })
  })
}

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <BrowserRouter><AuthProvider><App /></AuthProvider></BrowserRouter>
  </StrictMode>,
)
