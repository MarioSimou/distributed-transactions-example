import React from 'react'

export const initUserValues = {}

export const UserProfileContext = React.createContext(initUserValues)
export const useUserProfile = () => React.useContext(UserProfileContext)
