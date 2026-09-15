import { healthApi } from '@/lib/api/endpoints'

export default async function Page() {
  const data = await healthApi.health()

  return (
    <div>
      <h1>Ideagen</h1>
      <p>Environment: {data.environment}</p>
      <p>Status: {data.status}</p>
    </div>
  )
}
