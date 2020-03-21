const ACTIVE_NETWORK = 'active_network'

export const getNetwork = () => {
  var net = window.localStorage.getItem(ACTIVE_NETWORK)
  return JSON.parse(net)
}

export const saveNetwork = network => {
  window.localStorage.setItem(ACTIVE_NETWORK, JSON.stringify(network))
}

export const destroyNetwork = () => {
  window.localStorage.removeItem(ACTIVE_NETWORK)
}

export default { getNetwork, saveNetwork, destroyNetwork }
