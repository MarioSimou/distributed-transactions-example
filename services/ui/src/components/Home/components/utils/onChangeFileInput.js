import base64 from './base64'

const onChangeFileInput = appendField => async e => {
  try{
    const base64Image = await base64.encode(e.target.files[0])
    appendField('productImage', base64Image)
  }catch(e){
    console.warn(e.message)
  }
}

export default onChangeFileInput