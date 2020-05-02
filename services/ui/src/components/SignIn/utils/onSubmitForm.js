import httpClient from '../../../utils/httpClient.js'
import history from '../../../utils/history.js'

export default (formValues, setUserProfile) => async e => {
  e.preventDefault()

  try {
    const {data, status, message} = await httpClient({
      method: 'POST',
      url: new URL("/api/v1/signin", process.env.REACT_APP_CUSTOMERS_API),
      data: JSON.stringify({
        email: formValues.email.value,
        password: formValues.password.value,
      })
    })
    if (status !== 200){
      throw new Error(message)
    }

    setUserProfile(data.data)
    history.push('/')
  }catch(e){
    window.alert((e.response && e.response.data && e.response.data.message) || (e.message) )
  }
}