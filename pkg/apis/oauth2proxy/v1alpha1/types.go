package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Proxy `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Proxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ProxySpec   `json:"spec"`
	Status            ProxyStatus `json:"status,omitempty"`
}

type ProxySpec struct {
	Size   int32       `json:"size"`
	Config ProxyConfig `json:"config"`
	Image  string      `json:"image"`
}

type ProxyConfig struct {
	// CLI args
	CookieSecure string `json:"cookieSecure" cli:"--cookie-secure"`
	Upstream     string `json:"upstream" cli:"--upstream"`
	HTTPAddress  string `json:"address" cli:"--http-address"`
	RedirectURL  string `json:"redirectURL" cli:"--redirect-url"`
	EmailDomain  string `json:"emailDomain" cli:"--email-domain"`
	Provider     string `json:"provider" cli:"--provider"`
	// ENV variables
	CookieSecret string `json:"cookieSecret" env:"OAUTH2_PROXY_COOKIE_SECRET"`
	CookieDomain string `json:"cookieDomain" env:"OAUTH2_PROXY_COOKIE_DOMIAN"`
	ClientID     string `json:"clientID" env:"OAUTH2_PROXY_CLIENT_ID"`
	ClientSecret string `json:"clientSecret" env:"OAUTH2_PROXY_CLIENT_SECRET"`
	//github specific
	GithubOrg  string `json:"githubOrg" cli:"--github-org"`
	GithubTeam string `json:"githubString" cli:"--github-string"`
}

type ProxyStatus struct {
	// Fill me
}
