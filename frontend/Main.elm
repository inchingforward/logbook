module Main exposing (..)

import Html exposing (Html, button, div, input, text)
import Html.App as App
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import Http


main : Program Never
main =
    App.beginnerProgram { model = model, view = view, update = update }



-- Model


type alias Model =
    { username : String
    , password : String
    }


model : Model
model =
    { username = ""
    , password = ""
    }


type Msg
    = ChangeUsername String
    | ChangePassword String
    | Login
    | LoginSucceeded
    | LoginFailed Http.Error



-- Update


update : Msg -> Model -> Model
update msg model =
    case msg of
        ChangeUsername username ->
            { model | username = username }

        ChangePassword password ->
            { model | password = password }

        Login ->
            model

        LoginSucceeded ->
            model

        LoginFailed _ ->
            model



-- View


view : Model -> Html Msg
view model =
    div []
        [ text "Login"
        , input [ type' "text", placeholder "Username", onInput ChangeUsername ] []
        , input [ type' "password", placeholder "Password", onInput ChangePassword ] []
        , div [] [ text ("User: " ++ model.username ++ " Password: " ++ model.password) ]
        , button [ onClick Login ] [ text "Login" ]
        ]
