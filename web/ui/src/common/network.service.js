const ACTIVE_NETWORK = 'active_network'

export const getNetwork = () => {
  return window.localStorage.getItem(ACTIVE_NETWORK)
}

export const saveNetwork = network => {
  window.localStorage.setItem(ACTIVE_NETWORK, network)
}

export const destroyNetwork = () => {
  window.localStorage.removeItem(ACTIVE_NETWORK)
}

export default { getNetwork, saveNetwork, destroyNetwork }
