import httpClient from '../../../utils/httpClient.js'
import history from '../../../utils/history.js'

export default formValues => async e => {
  e.preventDefault()
  console.warn("DO SOME VALIDATION")

  try {
    const url = new URL('/api/v1/users', process.env.REACT_APP_CUSTOMERS_API).toString()
    const {data} = await httpClient({
      method: 'POST',
      url: url,
      data: JSON.stringify({
        username: formValues.username.value,
        email: formValues.email.value,
        password: formValues.password.value,
        confirmPassword: formValues.confirmPassword.value,
      })
    })

    if(!data.success) {
      throw new Error(data.message)
    }

    // store user data
    // console.log("USER:", data.data)
    history.push('/')

  } catch(e){
    window.alert( (e.response && e.response.data && e.response.data.message) || (e.message))
  }
}