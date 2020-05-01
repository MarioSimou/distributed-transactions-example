import React from 'react'
import {
  Paper,
  TextField,
  InputAdornment,
  Button,
  Typography,
} from '@material-ui/core'
import makeStyles from '@material-ui/styles/makeStyles'
import AccountCircle from '@material-ui/icons/AccountCircle';
import LockIcon from '@material-ui/icons/Lock';
import httpClient from '../utils/httpClient.js'
import history from '../utils/history'
import { useUserProfile } from '../utils/hooks'

const onSubmitForm = (formValues, setUserProfile) => async e => {
  e.preventDefault()

  try {
    const {data, headers, status, message} = await httpClient({
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
    window.alert(e.response && e.response.data && e.response.data.message || e.message )
  }
}

const CustomAdornment = ({position = "start", Icon}) => {
  return (
    <InputAdornment position={position}>
      <Icon/>
    </InputAdornment>
  )
}

const initFormValues = {
  email: {
    touched: false,
    value: ''
  },
  password: {
    touched: false,
    value: ''
  },
}

const Login = () => {
  const classes = useStyles()
  const [formValues, setFormValues] = React.useState(initFormValues) 
  const handleOnChangeField = key => ({target}) => setFormValues({...formValues, [key]: {touched: true, value: target.value}})
  const handleOnFocusField = key => () => setFormValues({...formValues, [key]: {touched: true, value: formValues[key]['value']}})
  const onChangeEmail = handleOnChangeField('email') 
  const onChangePassword = handleOnChangeField('password') 
  const onFocusEmail = handleOnFocusField('email')
  const onFocusPassword = handleOnFocusField('password')
  const {setUserProfile} = useUserProfile()

  return (
    <Typography component="div" className={classes.root}>
      <Paper elevation={3} variant="outlined" className={classes.paper}>
        <Typography variant="h3" align="center" className={classes.title}>Sign In</Typography>
        <form className={classes.form} noValidate={true} autoComplete="off" onSubmit={onSubmitForm(formValues, setUserProfile)}>
          <TextField id="email" 
                    type="email"
                    label="Email" 
                    placeholder="Your email" 
                    variant="filled" 
                    value={formValues.email.value}
                    error={formValues.email.touched && formValues.email.value === ''}
                    onChange={onChangeEmail}
                    onFocus={onFocusEmail}
                    InputProps={{startAdornment: <CustomAdornment Icon={AccountCircle}/>}} 
                    fullWidth
                    required/>
          <TextField id="password" 
                    type="password"
                    label="Password" 
                    placeholder="Your Password" 
                    variant="filled" 
                    error={formValues.password.touched && formValues.password.value === ''}
                    value={formValues.password.value}
                    onChange={onChangePassword}
                    onFocus={onFocusPassword}
                    InputProps={{startAdornment: <CustomAdornment Icon={LockIcon}/>}} 
                    fullWidth
                    required/>
          <Button type="submit" variant="contained" color="primary" fullWidth className={classes.button}>
            Login
          </Button>
        </form>
      </Paper>
    </Typography>
  )
}

const useStyles = makeStyles(theme => ({
  root: {
    minHeight: 'calc( 100vh - 48px )',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    background: theme.palette.primary.light,
  },
  paper: {
    padding: theme.spacing(3),
  },
  title: {
    marginBottom: theme.spacing(2),
  },
  form: {
    display: 'grid',
    gridAutoFlow: 'row',
    gridRowGap: theme.spacing(2),
    maxWidth: 500,
    width: 500,
    borderRadius: theme.spacing(1),
  },
  button: {
    marginTop: theme.spacing(2),
    height: 40
  }
}))

export default Login