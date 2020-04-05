import format from 'date-fns/format'

export default time => {
  return format(new Date(time), 'YYYY-MM-DD hh:mm:ss')
}
