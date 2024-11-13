const SERVER_INFO_KEY = "SERVER_INFO"

// TODO: Abstract state logic

/** @type State **/
let state = {
  serverList: []
}

/** @type Function[] **/
const listeners = []

/** @param {State} newState **/
const setState = (newState) => {
  Object.assign(state, newState);
  listeners.forEach((listener) => listener(state))
}

/** @param {StateListener} listener **/
const subscribe = (listener) => {
  listeners.push(listener)
  listener(state)
}

/**
  * Fetch server list
  *
  * @return {Promise<void>}
  *
  * **/
const fetchServerList = async () => {
  try {
    const response = await fetch("/api/servers")

    const serverList = await response.json()

    setState({
      serverList
    })
  } catch (error) {
    console.error("Error fetching server list")
  }
}

const initializeServerList = async () => {
  const table = document.getElementById("server-list")
  table.innerHTML = ""

  const thead = `<tr>
    <th class="border border-solid border-blue-800/70 lg:border-blue-800/40 p-2 text-center">Name</th>
    <th class="border border-solid border-blue-800/70 lg:border-blue-800/40 p-2 text-center">Players</th>
    <th class="border border-solid border-blue-800/70 lg:border-blue-800/40 p-2 text-center">IP Address:Port</th>
    <th class="border border-solid border-blue-800/70 lg:border-blue-800/40 p-2 text-center">Ping</th>
    <th class="border border-solid border-blue-800/70 lg:border-blue-800/40 p-2 text-center">Map</th>
    </tr>`

  const tableHead = document.createElement("thead")
  tableHead.innerHTML = thead

  table.appendChild(tableHead)

  subscribe(({ serverList }) => {
    serverList.forEach(({ ServerInfo, Ping }) => {
      const serverName = ServerInfo.Details.Name
      const players = `${ServerInfo.Details.Players}/${ServerInfo.Details.MaxPlayers}`
      const ipAddress = `${ServerInfo.IpAddress}`
      const ping = `${Ping}`
      const mapImage = ServerInfo.MapImage
      const map = ServerInfo.Details.Map

      const rowData = `
        <td class="border-[1px] border-solid border-blue-800/70 lg:border-blue-800/40 p-2 text-center">${serverName}</td>
        <td class="border-[1px] border-solid border-blue-800/70 lg:border-blue-800/40 p-2 text-center">${players}</td>
        <td class="border-[1px] border-solid border-blue-800/70 lg:border-blue-800/40 p-2 text-center">${ipAddress}</td>
        <td class="border-[1px] border-solid border-blue-800/70 lg:border-blue-800/40 p-2 text-center">${ping}</td>
        <td class="border-b border-blue-800/70 lg:border-blue-800/40 border-r flex flex-col p-2 gap-2">
        <img src="${mapImage}" alt="Map Image" class="mx-auto">
        <p class="text-center">${map}</p>
        </td>
        `

      const row = document.createElement("tr")
      row.innerHTML = rowData
      table.appendChild(row)
    })
  })
}

/** @param {RegisterServerResponse} serverInfo **/
const setServerDetails = (serverInfo) => {
  localStorage.setItem(SERVER_INFO_KEY, JSON.stringify(serverInfo))
}

const resetServerDetails = () => {
  localStorage.removeItem(SERVER_INFO_KEY)
}

/**
  *
  * @returns {RegisterServerResponse | undefined}
  *
  * **/
const getServerDetails = () => {
  const serverInfo = localStorage.getItem(SERVER_INFO_KEY)
  try {
    return /** @type RegisterServerResponse **/ (JSON.parse(serverInfo))
  } catch (error) {
    console.error("Error parsing serverInfo")
    return undefined
  }
}


/** @param {RegisterServerResponse} serverInfo **/
const populateSuccessSectionDom = (serverInfo) => {
  const copyIpInput = /** @type HTMLInputElement **/ (document.getElementById("copy-ip-input"))
  const copyNickInput = /** @type HTMLInputElement **/ (document.getElementById("copy-nick-input"))
  const copyPwInput = /** @type HTMLInputElement **/ (document.getElementById("copy-pw-input"))

  const getStartedSection = document.getElementById("get-started-section")
  const registerForm = document.getElementById("register-server-form")
  const successSection = document.getElementById("success-section")

  copyIpInput.value = serverInfo.IpAddress
  copyNickInput.value = serverInfo.AdminNickname
  copyPwInput.value = serverInfo.AdminPassword

  getStartedSection.style.display = "none"
  registerForm.style.display = "none"
  successSection.style.display = "flex"
}

const resetSuccessSectionDom = () => {
  const copyIpInput = /** @type HTMLInputElement **/ (document.getElementById("copy-ip-input"))
  const copyNickInput = /** @type HTMLInputElement **/ (document.getElementById("copy-nick-input"))
  const copyPwInput = /** @type HTMLInputElement **/ (document.getElementById("copy-pw-input"))

  const getStartedSection = document.getElementById("get-started-section")
  const registerForm = document.getElementById("register-server-form")
  const successSection = document.getElementById("success-section")

  copyIpInput.value = ""
  copyNickInput.value = ""
  copyPwInput.value = ""

  successSection.style.display = "none"
  registerForm.style.display = "flex"
  getStartedSection.style.display = "flex"
}

const initializeRegisterServerForm = () => {
  // Type assert the fool
  const registerServerForm = /** @type HTMLFormElement **/ (document.getElementById("register-server-form"))

  registerServerForm.onsubmit = (e) => {
    const onSubmit = async () => {
      const formData = new FormData(registerServerForm)

      const maxPlayers = /** @type string **/ (formData.get("maxPlayers"))
      const adminNickname = formData.get("adminNickname")

      const submitBtn = /** @type HTMLButtonElement **/ (document.getElementById("submit-btn"))
      const submitSvg = document.getElementById("submit-svg")

      const submitInfo = document.getElementById("submit-info")


      submitSvg.classList.add("animate-spin")
      submitBtn.innerHTML = "Processing..."
      submitBtn.disabled = true
      submitInfo.style.display = "block"

      try {
        const response = await fetch("/api/register", {
          method: "POST",
          body: JSON.stringify({ maxPlayers: Number.parseInt(maxPlayers), adminNickname })
        })

        /** @type RegisterServerResponse **/
        const json = await response.json()

        setServerDetails(json)
        populateSuccessSectionDom(json)
      } catch (error) {
        console.error("Error registering server")
      } finally {
        submitSvg.classList.remove("animate-spin")
        submitBtn.innerHTML = "Submit"
        submitBtn.disabled = false
        submitInfo.style.display = "none"


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

const initializeSavedServerDetails = () => {
  const serverDetails = getServerDetails()

  if (!serverDetails) {
    resetSuccessSectionDom()
    return
  }

  subscribe(({ serverList }) => {
    const hasProvisionedServerExpired = !serverList.some(({ ServerInfo }) => ServerInfo.IpAddress === serverDetails.IpAddress)

    if (hasProvisionedServerExpired) {
      resetServerDetails()
      resetSuccessSectionDom()
      return
    }

    populateSuccessSectionDom(serverDetails)
  })
}

initializeServerList()
initializeRegisterServerForm()
initializeCopyButtonsListeners()
initializeSavedServerDetails()
fetchServerList()
