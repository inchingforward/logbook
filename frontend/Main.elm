module Main exposing (..)

import Html exposing (Html, button, div, input, text)
import Html.App as App
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Decode exposing (..)
import Json.Encode as JSEncode
import Task


main : Program Never
main =
    App.program
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }



-- Model


type alias Model =
    { username : String
    , password : String
    , serverResult : String
    }


type alias LoginResult =
    { success : Bool
    , message : String
    , token : String
    }


model : Model
model =
    { username = ""
    , password = ""
    , serverResult = ""
    }


init : ( Model, Cmd Msg )
init =
    ( { username = "", password = "", serverResult = "" }, Cmd.none )



-- Update


type Msg
    = ChangeUsername String
    | ChangePassword String
    | Login
    | LoginSucceeded LoginResult
    | LoginFailed Http.Error


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        ChangeUsername username ->
            ( { model | username = username }, Cmd.none )

        ChangePassword password ->
            ( { model | password = password }, Cmd.none )

        Login ->
            ( model, login model )

        LoginSucceeded result ->
            ( { model | serverResult = "Logged in!" }, Cmd.none )

        LoginFailed error ->
            ( { model | serverResult = toString error }, Cmd.none )


encodeModel : Model -> JSEncode.Value
encodeModel model =
    let
        list =
            [ ( "username", JSEncode.string model.username )
            , ( "password", JSEncode.string model.password )
            ]
    in
        list |> JSEncode.object


decodeLoginResult : Decoder LoginResult
decodeLoginResult =
    object3 LoginResult
        ("success" := bool)
        ("message" := string)
        ("token" := string)


login : Model -> Cmd Msg
login model =
    Task.perform
        LoginFailed
        LoginSucceeded
        ((Http.send
            Http.defaultSettings
            { verb = "POST"
            , headers = [ ( "Content-Type", "application/json" ) ]
            , url = "http://localhost:4003/login"
            , body = model |> encodeModel |> JSEncode.encode 0 |> Http.string
            }
         )
            |> Http.fromJson decodeLoginResult
        )



-- Subscriptions


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- View


view : Model -> Html Msg
view model =
    div []
        [ text "Login"
        , input [ type' "text", placeholder "Username", onInput ChangeUsername ] []
        , input [ type' "password", placeholder "Password", onInput ChangePassword ] []
        , div [] [ text ("User: " ++ model.username ++ " Password: " ++ model.password) ]
        , div [] [ text model.serverResult ]
        , button [ onClick Login ] [ text "Login" ]
        ]
