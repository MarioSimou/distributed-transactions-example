import httpClient from '../../../utils/httpClient'
import history from '../../../utils/history'
import {initUserValues} from '../../../utils/hooks.js'

export default ({setUserProfile}) => async () => {
  try {
    const uri = new URL(`/api/v1/logout`, process.env.REACT_APP_CUSTOMERS_API)
    const {status} = await httpClient.post(uri.toString())
    if (status === 204){
      setUserProfile(initUserValues)
      history.push('/')
    }
  }catch(e){
    window.alert((e.response && e.response.data && e.response.data.message) || (e.message))
  }
}