module Main exposing (..)

import Html exposing (Html, button, div, input, text)
import Html.App as App
import Html.Attributes exposing (..)
import Html.Events exposing (onClick)


main =
    App.beginnerProgram { model = 0, view = view, update = update }


type Msg
    = Login


update msg model =
    model


view model =
    div []
        [ text "Login"
        , input [ type' "text", placeholder "Username" ] []
        , input [ type' "password", placeholder "Password" ] []
        , button [ onClick Login ] [ text "Login" ]
        ]
