import format from 'date-fns/format'

export default time => {
  return format(new Date(time), 'YYYY-MM-D hh:mm:ss')
}
