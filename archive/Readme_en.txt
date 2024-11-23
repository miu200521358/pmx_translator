----------------------------------------------------------------
----------------------------------------------------------------

�@�uPmxTranslator�v

�@�@ver1.00.00

�@�@�@�@�@�@�@�@�@�@�@�@�@�@�@�@�@miu200521358

----------------------------------------------------------------
----------------------------------------------------------------

Thank you for downloading my work.
Please take a moment to review the following before using it.


----------------------------------------------------------------


----------------------------------------------------------------
���@Summary
----------------------------------------------------------------

This is a tool that allows you to batch replace the name fields in PMX models.
It is primarily intended for converting Chinese model names into Japanese, but there are no restrictions on loading models, so it can also be used for renaming existing models.


----------------------------------------------------------------
���@Distribution Video
----------------------------------------------------------------

�iworking�j


----------------------------------------------------------------
���@Static Images for Content Tree
----------------------------------------------------------------

�iworking�j


----------------------------------------------------------------
���@Included Files
----------------------------------------------------------------

- PmxTranslator_x.xx.xx.exe�@�c�@exe file
- Readme.txt�@�@�@�@�@�@ �c�@Readme


----------------------------------------------------------------
���@System Requirements
----------------------------------------------------------------

- Windows10/11 64bit
  - This tool is intended for use in environments where MMD (MikuMikuDance) is supported.


----------------------------------------------------------------
���@How to use
----------------------------------------------------------------

- Basically, you can simply run the exe as is.
- File history can be copied by placing "user_config.json" in the same directory as the exe.
- Name Replacement
    - Specify the PMX model and CSV dictionary to replace names that match the specified strings.
    - The output path is also replaced according to the dictionary, and it outputs everything, including textures.
    - When you load the model and dictionary, the replacement results are displayed in the table field.
    - The background color of the name fields that are scheduled to be replaced is changed.
- CSV Output
    - Outputs CSV data, which serves as the source data for the dictionary.
    - When you load the PMX model, the list of names is displayed in the table field.
    - If there are characters that cannot be recognized by Shift-JIS (such as Simplified Chinese characters), a check is added when initially displayed.
    - You can also include any other name fields you want to convert by checking them.

For detailed usage instructions, please refer to the video.

----------------------------------------------------------------
���@How to create the dictionary
----------------------------------------------------------------

https://chatgpt.com/g/g-Cmchsxm6X-pmxtranslaordictionarygenerator

- You can use ChatGPT to complete the dictionary CSV.
- Open the URL above, input the text output from "CSV Export," and press the Enter key to get the completion results.
- Copy the completion results and paste them into the original CSV file to create a dictionary that can be used for "Name Replacement."

Note: You will need to create a ChatGPT account (free).

----------------------------------------------------------------
���@Community Information
----------------------------------------------------------------

MMD�����Fhttps://discord.gg/jDqU6qRaKj

�@This is a unique Discord server where each category has an owner, and you can join the category of your choice.
�@It is intended for use in distributing community-limited content and interacting, such as through the Nico Nico Community.
�@You can also apply to open new categories.

miu�̎������Fhttps://discord.gg/MW2Bn47aCN

�@This lab focuses on various experiments with self-made tools, such as VMD sizing and Motion Supporter.
�@I welcome questions and consultations regarding my tools. An FAQ page is also available.
�@Additionally, you can try out beta versions of tools in advance.
�@It is one category in the "MMD Nagaya" server, so please join the server and then apply to join this category.


----------------------------------------------------------------
���@Terms of Use and Other Information
----------------------------------------------------------------

�sMandatory Requirements�t

- If you publish or distribute a modified model, please include proper credit.
- For Nico Nico Douga, please register the work in the content tree using a tree image (currently in preparation).
  - If you register as a parent in the content tree, credit listing is optional.
- When distributing the model to an unspecified large number of people, please include credit or register it in the content tree only in the source (e.g., video) announcing the distribution.
  - It is not necessary to list credits for works using the relevant model.

�sOptional Items�t

- Regarding this tool and modified models, you are free to do the following within the scope of the original model's terms:

- Adjusting and modifying the modified model
  - For distributed models, please check if modifications are permitted by the original terms.
- Posting videos using the model on video-sharing sites, social media, etc.
  - If the original model's terms specify conditions such as posting destinations or age restrictions, the modified model created with this tool must also comply with these conditions.
- Distributing the modified model to an unspecified large number of people
  - Only for self-made models or models that are permitted to be distributed to an unspecified large number of people.

�sProhibited Items�t

- Regarding this tool and the models generated with it, please refrain from the following actions:

  - Actions beyond the scope of the original model's terms.
  - Claiming the work as entirely self-made.
  - Actions that may cause inconvenience to the original rights holders.
  - Using the models for the purpose of defaming or slandering others (regardless of whether it is two-dimensional or three-dimensional).
  
  - The following are not prohibited, but I ask for your consideration:
    - Using the models in works containing excessive violence, obscenity, romantic content, grotesque elements, political or religious expressions (equivalent to R-15 or above).
    - Please ensure that your use of the tool and models complies with the original model's terms.
    - When publishing your work, please take care to avoid search engine indexing, etc.

�@- Please note that while 'use for commercial purposes' is not prohibited for this tool, it is prohibited under the terms of PMXEditor.

�sDisclaimer�t

- Please use the tool at your own risk.
- The author assumes no responsibility for any issues arising from the use of this tool.


----------------------------------------------------------------
���@Source Code & Libraries
----------------------------------------------------------------

This tool is created using Go, and includes the following libraries:

�Ehttps://github.com/lxn/walk
�Ehttps://github.com/go-gl/gl
�Ehttps://github.com/go-gl/glfw/v3.3/glfw
�Ehttps://github.com/go-gl/mathgl
�Ehttps://github.com/nicksnyder/go-i18n/v2
�Ehttps://github.com/petar/GoLLRB
�Ehttps://github.com/ftrvxmtrx/tga

The source code is available on GitHub (MIT License).
However, the copyright is not waived.

https://github.com/miu200521358/pmx_translator

I created the icon myself based on a rough draft from ChatGPT.


----------------------------------------------------------------
���@Credit
----------------------------------------------------------------

�@Tool Name�F PmxTranslator
�@Author�F�@miu (miu200521358)

�@http://www.nicovideo.jp/user/2776342
�@Twitter: @miu200521358
�@Mail: garnet200521358@gmail.com


----------------------------------------------------------------
���@Revision
----------------------------------------------------------------


