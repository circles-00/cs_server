
export { }

declare global {
  type ServerInfoResponse = {
    ServerInfo: ServerInfo
    PlayerInfo: PlayerInfo
    Ping: number
  }

  type ServerInfo = {
    Details: Details
    MapImage: string
    IpAddress: string
  }

  type Details = {
    Protocol: number
    Name: string
    Map: string
    Folder: string
    Game: string
    AppID: number
    Players: number
    MaxPlayers: number
    Bots: number
    ServerType: number
    ServerOS: number
    Visibility: boolean
    VAC: boolean
    Version: string
    EDF: number
    ExtendedServerInfo: ExtendedServerInfo
  }

  type ExtendedServerInfo = {
    Port: number
    SteamID: number
    Keywords: string
    GameID: number
  }

  type PlayerInfo = {
    Count: number
    Players: any
  }

  type RegisterServerResponse = {
    IpAddress: string
    AdminNickname: string
    AdminPassword: string
  }
}
