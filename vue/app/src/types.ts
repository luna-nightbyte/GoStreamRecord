export type Client = {
  id: string
  name: string
  url: string
}

export type PortalSettings = {
  clients: Client[]
  lastSelectedClientId?: string
}