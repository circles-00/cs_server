/**
  * Fetch server list
  *
  * @return {Promise<ServerInfoResponse[]>}
  *
  * **/
const fetchServerList = async () => {
  try {
    const response = await fetch("/api/servers")

    const json = await response.json()

    return json
  } catch (error) {
    console.error("Error fetching server list")
  }
}

const initializeServerList = async () => {
  const serverList = await fetchServerList()

  const table = document.getElementById("server-list")
  table.innerHTML = ""

  const thead = `<tr>
    <th class="border border-solid border-blue-800/40 p-2 text-center">Name</th>
    <th class="border border-solid border-blue-800/40 p-2 text-center">Players</th>
    <th class="border border-solid border-blue-800/40 p-2 text-center">IP Address:Port</th>
    <th class="border border-solid border-blue-800/40 p-2 text-center">Ping</th>
    <th class="border border-solid border-blue-800/40 p-2 text-center">Map</th>
    </tr>`

  const tableHead = document.createElement("thead")
  tableHead.innerHTML = thead

  table.appendChild(tableHead)

  serverList.forEach(({ ServerInfo, Ping }) => {
    const serverName = ServerInfo.Details.Name
    const players = `${ServerInfo.Details.Players}/${ServerInfo.Details.MaxPlayers}`
    const ipAddress = `${ServerInfo.IpAddress}`
    const ping = `${Ping}`
    const mapImage = ServerInfo.MapImage
    const map = ServerInfo.Details.Map

    const rowData = `
      <td class="border-[1px] border-solid border-blue-800/40 p-2 text-center">${serverName}</td>
      <td class="border-[1px] border-solid border-blue-800/40 p-2 text-center">${players}</td>
      <td class="border-[1px] border-solid border-blue-800/40 p-2 text-center">${ipAddress}</td>
      <td class="border-[1px] border-solid border-blue-800/40 p-2 text-center">${ping}</td>
      <td class="border-b border-blue-800/40 border-r flex flex-col p-2 gap-2">
      <img src="${mapImage}" alt="Map Image" class="mx-auto">
      <p class="text-center">${map}</p>
      </td>
      `

    const row = document.createElement("tr")
    row.innerHTML = rowData
    table.appendChild(row)
  })

}

const initializeRegisterServerForm = () => {
  // Type assert the fool
  const registerServerForm = /** @type HTMLFormElement **/ (document.getElementById("register-server-form"))

  registerServerForm.onsubmit = (e) => {
    const onSubmit = async () => {
      const formData = new FormData(registerServerForm)

      const maxPlayers = /** @type string **/ (formData.get("maxPlayers"))
      const adminNickname = formData.get("adminNickname")

      const submitBtn = document.getElementById("submit-btn")
      const submitSvg = document.getElementById("submit-svg")

      const submitInfo = document.getElementById("submit-info")
      const registerForm = document.getElementById("register-server-form")
      const successSection = document.getElementById("success-section")

      const copyIpInput = /** @type HTMLInputElement **/ (document.getElementById("copy-ip-input"))
      const copyNickInput = /** @type HTMLInputElement **/ (document.getElementById("copy-nick-input"))
      const copyPwInput = /** @type HTMLInputElement **/ (document.getElementById("copy-pw-input"))

      submitSvg.classList.add("animate-spin")
      submitBtn.innerHTML = "Processing..."
      submitInfo.style.display = "block"

      try {
        const response = await fetch("/api/register", {
          method: "POST",
          body: JSON.stringify({ maxPlayers: Number.parseInt(maxPlayers), adminNickname })
        })

        /** @type RegisterServerResponse **/
        const json = await response.json()

        copyIpInput.value = json.IpAddress
        copyNickInput.value = json.AdminNickname
        copyPwInput.value = json.AdminPassword


      } catch (error) {
        console.error("Error registering server")
      } finally {
        submitSvg.classList.remove("animate-spin")
        submitBtn.innerHTML = "Submit"
        submitInfo.style.display = "none"

        registerForm.style.display = "none"
        successSection.style.display = "flex"

        initializeServerList()
      }
    }

    e.preventDefault()
    onSubmit()
  }
}

const initializeCopyButtonsListeners = () => {
  /** @param {string} id **/
  const onClick = (id) => {
    const btn = document.getElementById(id)
    const copyBtn = document.getElementById(`${id}-copy`)
    const checkBtn = document.getElementById(`${id}-check`)
    const input = /** @type HTMLInputElement **/ (document.getElementById(`${id}-input`))

    btn.onclick = () => {
      copyBtn.style.display = "none"
      checkBtn.style.display = "block"
      navigator.clipboard.writeText(input.value)

      setTimeout(() => {
        copyBtn.style.display = "block"
        checkBtn.style.display = "none"
      }, 3000)
    }


  }

  onClick("copy-ip")
  onClick("copy-nick")
  onClick("copy-pw")
}

initializeServerList()
initializeRegisterServerForm()
initializeCopyButtonsListeners()
