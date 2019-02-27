# AWBeta

## Video demo
[![](http://img.youtube.com/vi/tUlizjq5xOg/0.jpg)](http://www.youtube.com/watch?v=tUlizjq5xOg "AWBeta | Demo")

## Screens
![Demo Image ](https://github.com/msorins/AWBeta/blob/master/0.png?raw=true "Demo Image")

![Demo Image ](https://github.com/msorins/AWBeta/blob/master/1.jpg?raw=true "Demo Image")



# Idea
Everybody is ordering more and more items from the internet and keeping track of their location throughout delivery is paramount. Unfortunately not all couriering companies offer good enough auto tracking options (at least in Romania, this is the case).

So, AWBot is a messenger bot that is going to give you the current status of your package, the history of all changes and an option to subscribe to any status modifications (the bot giving you message when an update comes through).

# How does it work

By providing an AWB, the BackEnd is going to automatically infer the providing couriering firm and crawl it's website in order to retrieve the status. (or call their API if they have one)

I am using WitAI in order to map out different kinds of AWBs to couriering firms and to decode user intent, that can be:
* request status of AWB
* request all history
* subscribe to changes from now on

If the couriering firm name cannot be inferred from AWB, I am simply asking them to manually provide the name.

These are the entities defined in Wit:
* ![Demo Image ](https://github.com/msorins/AWBeta/blob/master/2.png?raw=true "Demo Image")

The architecture of the program is build with modularity in mind, so that at anytime the messaging platform of the bot can be changed, also there can be added support for new couriering firms.

The interface for the both types of entities being:
```
type IChat interface {
	HandleMessages(messageReceivedCallBack func(string, string) []string)
	SendMessage(userId string, msgs []string)
}
```

```
type ISolver interface {
	updateStatuses() SolverResponse

	GetStatuses() ([]string, SolverResponse)
	GetLastStatus() ([]string, SolverResponse)
	GetAwb() string
}

```


# Technologies used
* WitAI

### BackEnd:
* GO


> The Project was realised between 2nd and 3rd year of University (summer of 2018)