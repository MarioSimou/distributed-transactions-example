const encode = file => new Promise((resolve,reject) => {
  const reader = new FileReader()
  reader.readAsDataURL(file)
  reader.onload = () => resolve(reader.result)
  reader.onerror = () => reject({message: 'Failre to base64 image'})
})

export default {
  encode,
}