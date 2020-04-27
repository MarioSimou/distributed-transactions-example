import React from 'react'
import {
  Paper,
  TextField,
  InputAdornment,
  Typography,
  Button,
} from '@material-ui/core'
import makeStyles from '@material-ui/styles/makeStyles'
import AccountCircle from '@material-ui/icons/AccountCircle';
import LockIcon from '@material-ui/icons/Lock';
import PersonIcon from '@material-ui/icons/Person';

const CustomAdornment = ({position = "start", Icon}) => {
  return (
    <InputAdornment position={position}>
      <Icon/>
    </InputAdornment>
  )
}

const initFormValues = {
  username: {
    touched: false,
    value: ''
  },
  email: {
    touched: false,
    value: ''
  },
  password: {
    touched: false,
    value: ''
  },
  confirmPassword: {
    touched: false,
    value: ''
  },
}

const SignUp = () => {
  const classes = useStyles()
  const [formValues, setFormValues] = React.useState(initFormValues) 
  const handleOnChangeField = key => ({target}) => setFormValues({...formValues, [key]: {touched: true, value: target.value}})
  const handleOnFocusField = key => () => setFormValues({...formValues, [key]: {touched: true, value: formValues[key]['value']}})
  const onChangeUsername = handleOnChangeField('username') 
  const onChangeEmail = handleOnChangeField('email') 
  const onChangePassword = handleOnChangeField('password') 
  const onChangeConfirmPassword = handleOnChangeField('confirmPassword') 
  const onFocusUsername = handleOnFocusField('username')
  const onFocusEmail = handleOnFocusField('email')
  const onFocusPassword = handleOnFocusField('password')
  const onFocusConfirmPassword = handleOnFocusField('confirmPassword')

  const onSubmitForm = e => {
    e.preventDefault()
    console.log("SUBMITTING FORM")
  }


  return (
    <Typography component="div" className={classes.root}>
    <Paper elevation={3} variant="outlined" className={classes.paper}>
      <Typography variant="h3" align="center" className={classes.title}>Sign In</Typography>
      <form className={classes.form} noValidate={true} autoComplete="off" onSubmit={onSubmitForm}>
      <TextField id="username" 
                  label="Username" 
                  placeholder="Your Username" 
                  variant="filled" 
                  value={formValues.username.value}
                  error={formValues.username.touched && formValues.username.value === ''}
                  onChange={onChangeUsername}
                  onFocus={onFocusUsername}
                  InputProps={{startAdornment: <CustomAdornment Icon={PersonIcon}/>}} 
                  fullWidth
                  required/>
        <TextField id="email" 
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
        <TextField id="confirmPassword" 
                   label="Password" 
                   placeholder="Your Password" 
                   variant="filled" 
                   error={formValues.confirmPassword.touched && formValues.confirmPassword.value === ''}
                   value={formValues.confirmPassword.value}
                   onChange={onChangeConfirmPassword}
                   onFocus={onFocusConfirmPassword}
                   InputProps={{startAdornment: <CustomAdornment Icon={LockIcon}/>}} 
                   fullWidth
                   required/>

        <Button type="submit" variant="contained" color="primary" className={classes.button} fullWidth>
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

export default SignUp