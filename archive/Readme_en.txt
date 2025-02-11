----------------------------------------------------------------
----------------------------------------------------------------

  「PmxTranslator」

    ver1.00.00

                                  miu200521358

----------------------------------------------------------------
----------------------------------------------------------------

Thank you for downloading my work.
Please take a moment to review the following before using it.


----------------------------------------------------------------


----------------------------------------------------------------
■  Summary
----------------------------------------------------------------

This is a tool that allows you to batch replace the name fields in PMX models.
It is primarily intended for converting Chinese model names into Japanese, but there are no restrictions on loading models, so it can also be used for renaming existing models.
If the file path contains characters outside of Shift-JIS, the file name will be modified accordingly.


----------------------------------------------------------------
■  Distribution Video
----------------------------------------------------------------

【MMD】Pmxの名称欄を一括置換できるツールを作ってみた【PmxTranslator ver1.00】
https://www.nicovideo.jp/watch/sm44133638


----------------------------------------------------------------
■  Static Images for Content Tree
----------------------------------------------------------------

PmxTranslator - コンテンツツリー用静画
https://youtu.be/Ne0MqQDOzYs
https://seiga.nicovideo.jp/seiga/im11561836


----------------------------------------------------------------
■  Included Files
----------------------------------------------------------------

- PmxTranslator_x.xx.xx.exe     …  exe file
- Readme.txt                    …  Readme
- content_tree.txt              …  Content tree set.
- dict_csv_sample.csv           …  Sample dictionary Csv.
- Community Link - miu's Lab    … Invitation link to miu's Lab Discord server.
- Chinese-to-Japanese Conversion Dictionary - Provided by Kanna @MMD JC at their BOOTH store …  A dictionary for use with PmxTranslator shared by Kanna.


----------------------------------------------------------------
■  System Requirements
----------------------------------------------------------------

- Windows10/11 64bit
  - This tool is intended for use in environments where MMD (MikuMikuDance) is supported.


----------------------------------------------------------------
■  How to use
----------------------------------------------------------------

- Basically, you can run the .exe as-is without any additional setup.
- File history can be copied by placing user_config.json in the same directory as the .exe.

- Name Replacement
  - Specify a Pmx model and a Csv dictionary to replace names containing matching strings.
  - The output path is also replaced according to the dictionary, including textures and other related files.
  - When the model and dictionary are loaded, the replacement results are displayed in the table section.
  - The fields where replacement is planned have a changed background color.
  - Clicking on any row outside the checkbox will open a dialog to convert Japanese names and English names individually.
    The Japanese field cannot be left empty. (Leaving the English field blank will not cause an error.)
- Csv Output
  - Outputs Csv data as source data for the dictionary.
  - When a Pmx model is loaded, the name list is displayed in the table section.
  - If characters that cannot be recognized in Shift-JIS (e.g., simplified Chinese characters) are included, a checkmark is set during the initial display.
  - You can include other fields you want to convert by checking them as output targets.
- Csv Merge
  - Combines two Csv files and outputs a new Csv without duplicates.

For detailed instructions, please refer to the video.


----------------------------------------------------------------
■  How to create the dictionary
----------------------------------------------------------------

https://chatgpt.com/g/g-Cmchsxm6X-pmxtranslaordictionarygenerator

- You can use ChatGPT to complete the dictionary CSV.
- Open the URL above, input the text output from "CSV Export," and press the Enter key to get the completion results.
- Copy the completion results and paste them into the original CSV file to create a dictionary that can be used for "Name Replacement."

Note: You will need to create a ChatGPT account (free).

----------------------------------------------------------------
■  Community Information
----------------------------------------------------------------

MMD長屋：https://discord.gg/jDqU6qRaKj

  This is a unique Discord server where each category has an owner, and you can join the category of your choice.
  It is intended for use in distributing community-limited content and interacting, such as through the Nico Nico Community.
  You can also apply to open new categories.

miuの実験室：https://discord.gg/MW2Bn47aCN

  This lab focuses on various experiments with self-made tools, such as VMD sizing and Motion Supporter.
  I welcome questions and consultations regarding my tools. An FAQ page is also available.
  Additionally, you can try out beta versions of tools in advance.
  It is one category in the "MMD Nagaya" server, so please join the server and then apply to join this category.


----------------------------------------------------------------
■  Terms of Use and Other Information
----------------------------------------------------------------

《Mandatory Requirements》

Please make sure to perform the following actions regarding this tool and the modified model:

- If you publish or distribute a modified model, please include proper credit.
  - It would be helpful if it could be searched with "(Tool name or MMD) miu."
- For Nico Nico Douga, please register the work in the content tree using a tree image (currently in preparation).
  - If you register as a parent in the content tree, credit listing is optional.
- When distributing the model to an unspecified large number of people, please include credit or register it in the content tree only in the source (e.g., video) announcing the distribution.
  - It is not necessary to list credits for works using the relevant model.


《Optional Items》

Regarding this tool and modified models, you are free to do the following within the scope of the original model's terms:

- Adjusting and modifying the modified model
  - For distributed models, please check if modifications are permitted by the original terms.
- Posting videos using the model on video-sharing sites, social media, etc.
  - If the original model's terms specify conditions such as posting destinations or age restrictions, the modified model created with this tool must also comply with these conditions.
- Distributing the modified model to an unspecified large number of people
  - Only for self-made models or models that are permitted to be distributed to an unspecified large number of people.


《Prohibited Items》

Regarding this tool and the models generated with it, please refrain from the following actions:

  - Actions beyond the scope of the original model's terms.
  - Claiming the work as entirely self-made.
  - Actions that may cause inconvenience to the original rights holders.
  - Using the models for the purpose of defaming or slandering others (regardless of whether it is two-dimensional or three-dimensional).

  - The following are not prohibited, but I ask for your consideration:
    - Using the models in works containing excessive violence, obscenity, romantic content, grotesque elements, political or religious expressions (equivalent to R-15 or above).
    - Please ensure that your use of the tool and models complies with the original model's terms.
    - When publishing your work, please take care to avoid search engine indexing, etc.

  - Please note that while 'use for commercial purposes' is not prohibited for this tool, it is prohibited under the terms of PMXEditor.


《Disclaimer》

- Please use the tool at your own risk.
- The author assumes no responsibility for any issues arising from the use of this tool.


----------------------------------------------------------------
■  Source Code & Libraries
----------------------------------------------------------------

This tool is created using Go, and includes the following libraries:

・https://github.com/lxn/walk
・https://github.com/go-gl/gl
・https://github.com/go-gl/glfw/v3.3/glfw
・https://github.com/go-gl/mathgl
・https://github.com/nicksnyder/go-i18n/v2
・https://github.com/petar/GoLLRB
・https://github.com/ftrvxmtrx/tga

The source code is available on GitHub (MIT License).
However, the copyright is not waived.

https://github.com/miu200521358/pmx_translator

I created the icon myself based on a rough draft from ChatGPT.


----------------------------------------------------------------
■  Credit
----------------------------------------------------------------

  Tool Name： PmxTranslator
  Author：  miu (miu200521358)

  http://www.nicovideo.jp/user/2776342
  Twitter: @miu200521358
  Mail: garnet200521358@gmail.com


----------------------------------------------------------------
■  Revision
----------------------------------------------------------------

PmxTranslator_1.0.0 (2025/02/09)
- General distribution start.


