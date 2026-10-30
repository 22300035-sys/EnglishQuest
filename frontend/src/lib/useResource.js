import { useCallback, useEffect, useState } from 'react'
import { api, errorMessage } from './api.js'

export default function useResource(url) {
  const [state, setState] = useState({ data: null, loading: true, error: '' })
  const [version, setVersion] = useState(0)
  useEffect(() => {
    const controller = new AbortController()
    api
      .get(url, { signal: controller.signal })
      .then(({ data }) => {
        if (!controller.signal.aborted)
          setState({ data: data.data, loading: false, error: '' })
      })
      .catch((error) => {
        if (!controller.signal.aborted)
          setState({ data: null, loading: false, error: errorMessage(error) })
      })
    return () => controller.abort()
  }, [url, version])
  const reload = useCallback(() => {
    setState({ data: null, loading: true, error: '' })
    setVersion((n) => n + 1)
  }, [])
  return { ...state, reload }
}
