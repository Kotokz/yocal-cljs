(ns yocal-cljs.view.login
  (:require [reagent.core :as ratom]
            [clojure.string :as str]
            [re-frame.core :as re-frame]
            [reagent-forms.core :refer [bind-fields init-field value-of]]))

(def login-form-data (ratom/atom {:user {:username ""
                                         :password ""}
                                  :is-loading false}))

(defn row [label input]
  [:div.row
   [:div.col-md-2 [:label label]]
   [:div.col-md-5 input]])

(defn input [label type id]
   (row label [:input.form-control {:field type :id id}]))

(defn sign-in! []
  (let [name (get-in @login-form-data [:user :username])
        pwd (get-in @login-form-data [:user :password])]
    (cond
      (empty? name) (swap! login-form-data assoc-in [:errors :username] "user name is empty")
      (empty? pwd) (swap! login-form-data assoc-in [:errors :password] "password is empty")
      :else (re-frame/dispatch [:login (:user @login-form-data)]))))

(def form-login
  [:div
   (input "User Name" :text :user.username)
   [:div.row
    [:div.col-md-2]
    [:div.col-md-5
     [:div.alert.alert-danger
      {:field :alert :id :errors.username}]]]

   (input "password" :password :user.password)
   [:div.row
    [:div.col-md-2]
    [:div.col-md-5
     [:div.alert.alert-danger
      {:field :alert :id :errors.password}]
     [:div.alert.alert-danger
      {:field :alert :id :errors.other}]]]])


(def form-page
  (fn []
    (let [isloading (:is-loading @login-form-data)
          jwt (:jwt @(re-frame/subscribe [:user]))]
      (if-not (str/blank? jwt) (re-frame/dispatch [:set-active-panel :home-panel]) (pr "blank jwt"))
      [:div.col-md-6
       [:div.padding]
       [:div.page-header [:h1 "Login Form"]]
       [bind-fields
        form-login
        login-form-data]
       [:button.btn.btn-default {:on-click #(sign-in!)}
        (if isloading [:span.glyphicon.glyphicon-refresh.glyphicon-refresh-animate])
        "Submit"]])))
